package beater

import (
        "fmt"
        "net/url"
        "net/http"
	"crypto/tls"
        "io/ioutil"
        "encoding/json"
        "github.com/elastic/beats/libbeat/logp"
        jbcommon "github.com/kussj/jolokiabeat/common"
)

func (jb *Jolokiabeat) GetJMXMetrics(u string, qc []jbcommon.QueryConfig) (map[string]float64, error) {
        metrics := make(map[string]float64)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

        client := &http.Client{Transport: tr}

        /**
        Need to:
        x iterate over {domain, attributes} in qc
        x form rest calls to endpoint 'u + /read/ + domain + / + attr'
        - build map from responses form of: [domain + "/" + attribute => value]
        - return map
        */

        logp.Debug("Base URL: %s\n", u)

        for _,q := range qc {
                domain := q.GetDomain();
                attributes := q.GetAttributes()
                for _,attr := range attributes {
                        //var Url *url.URL
                        //Url,err := url.Parse(u)
                        //Url := url.QueryEscape(u + "/read/" + domain + "/" + attr)
                        //if err != nil {
                          //      return metrics, err
                        //}
                        //Url.Path += ("/read/" + domain + "/" + attr)
                        //logp.Debug("Encoded URL is %s\n", url.QueryEscape(u + "/read/" + domain + "/" + attr))
                        fmt.Println("> URL:", url.QueryEscape(u + "/read/" + domain + "/" + attr))
                }
        }


        resp, err := client.Get(u)
        defer resp.Body.Close()

        if err != nil {
                logp.Err("An error occured while executing HTTP request: %v", err)
                return metrics, err
        }

        // read json http response
	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)

        if err != nil {
                logp.Err("An error occured while reading HTTP response: %v", err)
                return metrics, err
        }

	err = json.Unmarshal([]byte(jsonDataFromHttp), &metrics)

        if err != nil {
                logp.Err("An error occured while unmarshaling jolokia data: %v", err)
                return metrics, err
        }

        return metrics, nil
}
