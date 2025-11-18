package handler

func (h *BaseHandler) StaticDataPath() string {
	return h.options.DataPath()
}
