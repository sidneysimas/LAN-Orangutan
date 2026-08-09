package web

import "html/template"

// typeIcons maps a device type to a small inline SVG, in the same Lucide style
// (fill none, currentColor stroke) as the other icons in the UI. The strings are
// static, never user input, so rendering them as template.HTML is safe.
var typeIcons = map[string]string{
	"Phone":        `<rect width="14" height="20" x="5" y="2" rx="2"/><path d="M12 18h.01"/>`,
	"Computer":     `<rect width="20" height="14" x="2" y="3" rx="2"/><line x1="8" x2="16" y1="21" y2="21"/><line x1="12" x2="12" y1="17" y2="21"/>`,
	"Printer":      `<path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/>`,
	"TV":           `<rect width="20" height="15" x="2" y="7" rx="2"/><polyline points="17 2 12 7 7 2"/>`,
	"Media Player": `<circle cx="12" cy="12" r="10"/><polygon points="10 8 16 12 10 16 10 8"/>`,
	"Router":       `<rect width="20" height="8" x="2" y="14" rx="2"/><path d="M6.01 18H6"/><path d="M10.01 18H10"/><path d="M15 10v4"/><path d="M17.84 7.17a4 4 0 0 0-5.66 0"/><path d="M20.66 4.34a8 8 0 0 0-11.31 0"/>`,
	"Server":       `<rect width="20" height="8" x="2" y="2" rx="2"/><rect width="20" height="8" x="2" y="14" rx="2"/><line x1="6" x2="6.01" y1="6" y2="6"/><line x1="6" x2="6.01" y1="18" y2="18"/>`,
	"IoT":          `<rect width="16" height="16" x="4" y="4" rx="2"/><rect width="6" height="6" x="9" y="9" rx="1"/><path d="M15 2v2M15 20v2M2 15h2M2 9h2M20 15h2M20 9h2M9 2v2M9 20v2"/>`,
	"Game Console": `<line x1="6" x2="10" y1="12" y2="12"/><line x1="8" x2="8" y1="10" y2="14"/><line x1="15" x2="15.01" y1="13" y2="13"/><line x1="18" x2="18.01" y1="11" y2="11"/><rect width="20" height="12" x="2" y="6" rx="2"/>`,
}

// typeIcon returns the inline SVG for a device type, or empty when the type is
// unknown. Used as a template function so the device list can show an icon
// beside each type badge.
func typeIcon(deviceType string) template.HTML {
	body, ok := typeIcons[deviceType]
	if !ok {
		return ""
	}
	return template.HTML(`<svg class="type-icon" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">` + body + `</svg>`)
}
