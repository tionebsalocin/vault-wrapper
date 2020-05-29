package runner

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/spnego"
)

func getKrb5Config(wrapper_config Config) string {
	content, err := ioutil.ReadFile(wrapper_config.Krb5)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func getToken(wrapper_config Config) string {
	l := log.New(os.Stderr, "GOKRB5 Client: ", log.LstdFlags)

	//defer profile.Start(profile.TraceProfile).Stop()
	// Load the keytab
	kt, err := keytab.Load(wrapper_config.Keytab)

	// Load the client krb5 config
	conf, err := config.NewFromString(getKrb5Config(wrapper_config))
	if err != nil {
		l.Fatalf("could not load krb5.conf: %v", err)
	}

	// Create the client with the keytab
	cl := client.NewWithKeytab(wrapper_config.VaultUser, wrapper_config.VaultUserRealm, kt, conf, client.Logger(l), client.DisablePAFXFAST(true))

	// Log in the client
	err = cl.Login()
	if err != nil {
		l.Fatalf("could not login client: %v", err)
	}

	// Form the request
	url := wrapper_config.VaultAddr + wrapper_config.VaultAuthPath
	r, err := http.NewRequest("POST", url, nil)
	if err != nil {
		l.Fatalf("could create request: %v", err)
	}

	spnegoCl := spnego.NewClient(cl, nil, wrapper_config.VaultSPN)

	// Make the request
	resp, err := spnegoCl.Do(r)
	if err != nil {
		l.Fatalf("error making request: %v", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Fatalf("error reading response body: %v", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		l.Fatalf("error decoding response body: %v", err)
	}

	auth := response["auth"].(map[string]interface{})
	return auth["client_token"].(string)
}
