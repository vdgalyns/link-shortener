package storage

var items = map[string]string{}

func Add(id string, value string) {
	if items == nil {
		items = map[string]string{}
	}
	items[id] = value
}

func Get(id string) (v string, ok bool) {
	if items == nil {
		items = map[string]string{}
	}
	v, ok = items[id]
	return
}
