package constants

var CommentOperators = map[string]string{
	".py":    "#",    // Python
	".go":    "//",   // Go
	".c":     "//",   // C
	".cpp":   "//",   // C++
	".h":     "//",   // C/C++ Header
	".cs":    "//",   // C#
	".java":  "//",   // Java
	".js":    "//",   // JavaScript
	".ts":    "//",   // TypeScript
	".php":   "//",   // PHP
	".rb":    "#",    // Ruby
	".rs":    "//",   // Rust
	".swift": "//",   // Swift
	".sh":    "#",    // Shell Script
	".pl":    "#",    // Perl
	".lua":   "--",   // Lua
	".m":     "%",    // MATLAB
	".r":     "#",    // R
	".scala": "//",   // Scala
	".kts":   "//",   // Kotlin
	".vb":    "'",    // Visual Basic .NET
	".f":     "!",    // Fortran
	".asm":   ";",    // Assembly
	".html":  "<!--", // HTML (Opening Comment)
	".css":   "/*",   // CSS (Opening Comment)
}
