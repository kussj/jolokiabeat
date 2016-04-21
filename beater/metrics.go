package beater

import (
        url "net/url"
        http "net/http"
	tls "crypto/tls"
        json "encoding/json"
        logp "github.com/elastic/beats/libbeat/logp"
        jbcommon "github.com/kussj/jolokiabeat/common"
)

type Request struct {
        Attribute       string `json:"attribute"`
        Mbean           string `json:"mbean"`
}

type Response struct {
        Status          int `json:"status"`
        Value           interface{} `json:"value"`
        Req             Request `json:"request"`
}

func (jb *Jolokiabeat) GetJMXMetrics(u string, qc []jbcommon.QueryConfig) (map[string]interface{}, error) {
        metrics := make(map[string]interface{})

//      logp.Debug("Base URL: %s\n", u)
        for _,q := range qc {
                domain := q.GetDomain();
                attributes := q.GetAttributes()
                for _,attr := range attributes {
                        var Url *url.URL
                        Url, err := url.Parse(u)
                        if err != nil {
                                logp.Err("An error occurred while parsing base url: %v", err)
                        }

                        Url.Path += ("/read/" + domain + "/" + attr)
                        //Url := url.QueryEscape(u + "/read/" + domain + "/" + attr)
                        //logp.Debug("> URL: ", Url.String())

                        resp := Response{};
                        err = getJson(Url.String(), &resp)

                        if err != nil {
                                logp.Err("An error occured while reading jolokia data: %v", err)
                                //return metrics, err
                        } else if resp.Status == 200 {
                                r := resp.Req
                                mbean := r.Mbean
                                attribute := r.Attribute
                                key := (mbean + "/" + attribute)
                                v := resp.Value
                                switch v.(type) {
                                case string:
                                        val := v.(string)
                                        metrics[key] = val
                                case float64:
                                        val := v.(float64)
                                        metrics[key] = val
                                case map[string]interface{}:
                                        val := mapWalker(v.(map[string]interface{}))
                                        switch val.(type){
                                        case string:
                                                val := val.(string)
                                                metrics[key] = val
                                        case float64:
                                                val := val.(float64)
                                                metrics[key] = val
                                        default:
                                                logp.Info(key, "is of a type I don't know how to handle")
                                        }//switch
                                default:
                                        logp.Info(key, "is of a type I don't know how to handle")
                                }//switch
                        }//else if
                }//for
        }//for

        return metrics, nil
}//GetJMXMetrics

func getJson(url string, target interface{}) error {
        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }//Transport

        client := &http.Client{Transport: tr}
        r, err := client.Get(url)
        if err != nil {
                logp.Err("An error occured while executing HTTP request: %v", err)
                return err
        }//if

        defer r.Body.Close()
        return json.NewDecoder(r.Body).Decode(target)
}//getJson

// Recursively parses through elements of a map until a 
// primative is found, then returns
func mapWalker(m map[string]interface{}) interface{} {
        for _,b := range m {
                switch b.(type) {
                case string:
                        return b.(string)
                case float64:
                        return b.(float64)
                case map[string]interface{}:
                        return mapWalker(b.(map[string]interface{}))
                }//switch
        }//for
        return nil
}//mapWalker

