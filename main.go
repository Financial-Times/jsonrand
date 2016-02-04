package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/jawher/mow.cli"
	"github.com/satori/go.uuid"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"time"
)

func main() {
	app := cli.App("jsonrand", "A randomised streaming json generator")
	template := app.StringOpt("template", "", "json example")
	count := app.IntOpt("count", 1, "how many json documents to generate")

	rand.Seed(time.Now().UnixNano())

	app.Action = func() {
		err := jsonrand(*template, *count)
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Run(os.Args)
}

func jsonrand(template string, count int) error {
	f, err := os.Open(template)
	if err != nil {
		return err
	}
	var j map[string]interface{}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&j); err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)

	for i := 0; i < count; i++ {
		randomizeMap(j)
		if err := enc.Encode(j); err != nil {
			return err
		}
	}

	return nil
}

func randomizeMap(j map[string]interface{}) map[string]interface{} {
	for k, v := range j {
		j[k] = randomizeValue(v)
	}
	return j
}

func randomizeValue(v interface{}) interface{} {
	switch x := v.(type) {
	case string:
		return randomizeString(x)
	case float64:
		return randomizeNumber(x)
	case map[string]interface{}:
		return randomizeMap(x)
	case []interface{}:
		for i, v := range x {
			x[i] = randomizeValue(v)
		}
		return x
	default:
		log.Panicf("bug. unhandled json type:%T", x)
		return nil
	}
}

var (
	uuidRegex = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
)

func randomizeNumber(f float64) float64 {
	ret := rand.Float64() * 9999
	// an integer?
	if math.Floor(f) == f {
		ret = math.Floor(ret)
	}
	return ret
}

func randomizeString(s string) string {

	// does this look like a time?
	if _, err := time.Parse(time.RFC3339, s); err == nil {
		r := time.Now().Unix() - time.Unix(0, 0).Unix()
		rnd := rand.Int63n(r)

		t := time.Unix(rnd, 0)
		return t.Format(time.RFC3339Nano)
	}

	// does this look like a uuid?
	if uuidRegex.MatchString(s) {
		return uuid.NewV4().String()
	}

	// otherwise fill with randomy stuff
	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := enc.Write(uuid.NewV4().Bytes()); err != nil {
		log.Panic(err)
	}
	return buf.String()
}
