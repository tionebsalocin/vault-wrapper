package runner

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func signalWatch(cmd *exec.Cmd) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		syscall.Kill(cmd.Process.Pid, sig.(syscall.Signal))
	}()
}

func getSecretsPathFromEnv(prefix string) []Secret {
	var secrets []Secret
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		if strings.Index(variable[0], prefix) == 0 && strings.Index(variable[1], "vault:") == 0 {
			secret_path := strings.Split(variable[1], ":")
			secrets = append(secrets, Secret{
				Name: strings.Replace(variable[0], prefix, "", 1),
				Path: secret_path[1],
				Key:  secret_path[2],
			})
		}
	}
	return secrets
}

func Run(config Config) {
	var secrets []string
	var secret string
	log.Println("[VaultWrapper] Authenticating to Vault...")
	token := getToken(config)
	if config.TokenOnly {
		os.Stdout.WriteString(token)
		return
	}
	paths := getSecretsPathFromEnv(config.EnvPrefix)
	if len(paths) > 0 {
		log.Println("[VaultWrapper] Collecting secrets...")
		vaultClient := vaultClient(config, token)
		for _, path := range paths {
			secret = getSecret(vaultClient, path.Path, path.Key)
			secrets = append(secrets, path.Name+"="+secret)
			if config.ExportOnly {
				os.Stdout.WriteString("export " + path.Name + "=" + secret + "\n")
			}
		}
	} else {
		log.Println("[VaultWrapper] No secret to collect")
	}
	if config.ExportOnly {
		os.Stdout.WriteString("export VAULT_TOKEN=" + token + "\n")
		return
	}
	cmd := exec.Command(config.CommandLine[0], config.CommandLine[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if len(secrets) > 0 {
		cmd.Env = os.Environ()
		//cmd.Env = []string{"coucou=pouet"}
		cmd.Env = append(cmd.Env, secrets...)
	}
	log.Println("[VaultWrapper] Launching command: ", config.CommandLine)
	if err := cmd.Start(); err != nil {
		log.Println("[VaultWrapper] Failed to start: ", err)
		return
	}
	log.Println("[VaultWrapper] Running with PID: ", cmd.Process.Pid)
	signalWatch(cmd)
	if err := cmd.Wait(); err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		log.Println("[VaultWrapper] Error: ", err.Error())
		log.Println("[VaultWrapper] Process stopped running. Exit code: ", exitCode)
		os.Exit(exitCode)
	}
	log.Println("[VaultWrapper] Process exited. Exit code: ", cmd.ProcessState.ExitCode())
}
