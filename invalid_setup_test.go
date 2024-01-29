package grom

func (c *Context) InvalidHandler() {}

func (c *Context) InvalidHandler2(w ResponseWriter, r *Request) string {
	return ""
}

func (c *Context) InvalidHandler3(w ResponseWriter, r ResponseWriter) {}
