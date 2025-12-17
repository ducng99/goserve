package pretty

//go:generate npx -y tailwindcss@3 -c ./tailwind.config.js -o ./tailwind.ignore.css --minify
//go:generate go run ../../../../../scripts/css_to_go.go tailwind.ignore.css
