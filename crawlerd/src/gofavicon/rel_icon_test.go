package gofavicon

import (
	"testing"
)

func TestRelIcon(t *testing.T) {
	i, _ := NewRelIcon("/img/favicon.png")
	if i.IsAbsURL() {
		t.Errorf("IsAbsURL failed for %v", i)
	}
	if i.IsEmbedded() {
		t.Errorf("IsEmbedded failed for %v", i)
	}

	i, _ = NewRelIcon("https://amazon.com/img/favicon.png")
	if !i.IsAbsURL() {
		t.Errorf("IsAbsURL failed for %v", i)
	}
}

func TestEmbeddedIcon(t *testing.T) {
	dataUrl := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyJpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMC1jMDYxIDY0LjE0MDk0OSwgMjAxMC8xMi8wNy0xMDo1NzowMSAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNS4xIFdpbmRvd3MiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6MzRFNEQ3NkJEMjRGMTFFMTk4RjA4NDhFNTEwRTcyMkIiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6MzRFNEQ3NkNEMjRGMTFFMTk4RjA4NDhFNTEwRTcyMkIiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDozNEU0RDc2OUQyNEYxMUUxOThGMDg0OEU1MTBFNzIyQiIgc3RSZWY6ZG9jdW1lbnRJRD0ieG1wLmRpZDozNEU0RDc2QUQyNEYxMUUxOThGMDg0OEU1MTBFNzIyQiIvPiA8L3JkZjpEZXNjcmlwdGlvbj4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gPD94cGFja2V0IGVuZD0iciI/PukMAkAAAAEsUExURcbf8uXx+vb6/eXx+Q5FtN/t+Ofz+rbX7uTr9+Tq9y5fvlh+y7rY7+n0+lJ6yS1evt7l9Sxdvr/b8L3a8BRLtpev31mAy+v0+k13yLLD59Hm9PT5/c3Z8LXG6Oby+pXG5aG54vX3/Git2LvZ79bg8vr8/vP4/Mjh8+Pw+Ofy+sPd8cfj9G2P0cDc783k9Ex2x6DK5xlOuCBUumWJztDm9DZkwLTF6Ojt92yP0s7a8OPq9m6Y1d3l9cHP7OLp9kNvxNPo9Wes2DNhwJSv3uDt+Ovw+cDP6yNVuurv+Nbf8pjI5aDL56G44rTL6h1RuWaJz9bp9tDo9u31+oWi2TZlwMfg8E16yd7t99Hc8Vl+ysri897l9JCq3FB5yPL1++/2/MLQ7HGb1j9sw////2AvJ1QAAADsSURBVHjaRM9ZMwNREAXgnpvLZKOENsQSJrTlIjHEEhQZS2JJxljHmoT+//9B37w4j19Xna4DzKzGDbHKLN+MdZlB4GNJ8fV9AbH1W7ew7g6zn9ORt/uM8wJkBOs7eHai/P05Bnp0S8ynBQz2qk9v0jEbTyjpKWZRI0Z9hqlV29u9ynQmR6PAY7g0NYH+efvzYDHAQ4bExBY0Bi3E7BHDkGkI+MXN7+ncV1pKU1AhIiYbudgvx957c+t2YW2k1ugJUKmqdXs7DGdiN5+yW+hV4K4cVhI37wzGdX4u0skL9FYcRw2ANh74P38CDABMCjqJfzUcfQAAAABJRU5ErkJggg=="
	i, _ := NewRelIcon(dataUrl)
	if !i.IsEmbedded() {
		t.Errorf("IsEmbedded failed for %v", i)
	}

	mimeType, bytes, err := i.Embedded()
	if err != nil {
		t.Error(err)
	}
	if mimeType != "image/png" {
		t.Errorf("Invalid mime type: %s", mimeType)
	}
	if len(bytes) == 0 {
		t.Errorf("Invalid bytes length")
	}
}

func TestEmbeddedExt(t *testing.T) {
	var expected map[string]string = map[string]string {
		"http://amazon.com/favicon.ico": ".ico",
		"/favicon.ico": ".ico",
		"https://website.com/sites/default/ico.png": ".png",
		"data:image/png;base64,aaaaa": "",
	}

	for u, e := range expected {
		r, _ := NewRelIcon(u)
		if ext := r.Ext(); ext != e {
			t.Errorf("invalid extension got \"%s\" expected \"%s\"", ext, e)
		}
	}

}