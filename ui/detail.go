package ui

import "github.com/kosuke9809/yadokari/sandbox"

type detailModel struct {
	sandbox sandbox.Sandbox
	showRaw bool
}

func newDetailModel() detailModel { return detailModel{} }

func (m detailModel) setSandbox(s sandbox.Sandbox) detailModel {
	m.sandbox = s
	return m
}

func (m detailModel) toggleRaw() detailModel {
	m.showRaw = !m.showRaw
	return m
}

func (m detailModel) view(_, _ int) string {
	if m.sandbox.ID == "" {
		return "No sandbox selected\n"
	}
	return "SANDBOX DETAIL: " + m.sandbox.Name + "\n"
}
