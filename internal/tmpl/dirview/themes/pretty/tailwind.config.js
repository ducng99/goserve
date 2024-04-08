/** @type {import('tailwindcss').Config} */
module.exports = {
    darkMode: 'selector',
    content: ['./*.templ'],
    theme: {
        extend: {
            gridTemplateColumns: {
                '4-1-1': '4fr 1fr 1fr',
            },
        },
    },
    plugins: [],
}

