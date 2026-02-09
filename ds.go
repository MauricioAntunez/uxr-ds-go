// Package ds provides a shared design system for UXR Go web applications.
//
// It includes reusable HTML component templates, CSS tokens and styles,
// vanilla JS utilities, and a unified template.FuncMap. Both uxr-backoffice
// and scrapibara/gobff consume this module to eliminate duplication.
//
// Usage:
//
//	import "uxr-ds"
//
//	// Get component templates filesystem
//	componentsFS := ds.Components()
//
//	// Get CSS/JS filesystems for static file serving
//	cssFS := ds.CSS()
//	jsFS := ds.JS()
//
//	// Get shared template functions
//	funcMap := ds.FuncMap()
//
//	// Merge with project-specific functions
//	funcMap = ds.MergeFuncMap(funcMap, myProjectFuncs)
package ds
