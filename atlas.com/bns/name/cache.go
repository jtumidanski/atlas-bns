package name

import (
	"log"
	"os"
	"sync"
)

type cache struct {
	names []string
}

func (c cache) GetNames() []string {
	var vs = make([]string, 0)
	for _, v := range c.names {
		vs = append(vs, v)
	}
	return vs
}

func (c cache) GetName(name string) (string, bool) {
	for _, v := range c.names {
		if v == name {
			return name, true
		}
	}
	return "", false
}

var c *cache
var once sync.Once

func InitCache(l *log.Logger) {
	once.Do(func() {
		d, e := os.LookupEnv("NAME_DIR")
		if !e {
			d = "/data/names"
		}
		vs, err := readDataDirectory(l, d)
		if err != nil {
			l.Fatal(err.Error())
		}

		l.Printf("[INFO] read %d companions.", len(vs))

		names := make([]string, 0)
		for _, v := range vs {
			l.Printf("[INFO] loading %s into cache.", v)
			names = append(names, v)
		}
		c = &cache{names: names}
	})
}

func GetCache() *cache {
	return c
}
