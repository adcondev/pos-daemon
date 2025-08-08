module.exports = {
    // Mapea los tipos a secciones del changelog
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "perf", section: "⚡ Performance Improvements"},
        {type: "refactor", section: "♻️ Code Refactoring"},
        {type: "docs", section: "📚 Documentation", hidden: true},
        {type: "test", section: "✅ Tests"},
        {type: "style", section: "💄 Style Changes", hidden: true},
        {type: "chore", section: "🔧 Maintenance", hidden: true},
        {type: "build", section: "🏗️ Build System"},
        {type: "ci", section: "👷 CI", hidden: true},
        {type: "deps", section: "📦 Dependencies"},
        {type: "revert", section: "⏪ Reverts"}
    ],

    // URLs del repositorio
    commitUrlFormat: "https://github.com/AdConDev/pos-daemon/commit/{{hash}}",
    compareUrlFormat: "https://github.com/AdConDev/pos-daemon/compare/{{previousTag}}...{{currentTag}}",
    issueUrlFormat: "https://github.com/AdConDev/pos-daemon/issues/{{id}}",
    userUrlFormat: "https://github.com/{{user}}",

    // Mensaje del commit de release
    releaseCommitMessageFormat: "chore(release): {{currentTag}}",

    // Prevenir saltos de versión accidentales
    skip: {
        bump: false,
        changelog: false,
        commit: false,
        tag: false
    },

    // Configuración de preset y parser
    preset: "conventionalcommits",
    presetConfig: {
        types: [
            {type: "feat", section: "✨ Features"},
            {type: "fix", section: "🐛 Bug Fixes"}
        ]
    },

    // Header del changelog
    header: "# Changelog\n\nAll notable changes to this project will be documented in this file.\n",

    // Configurar detección de breaking changes
    noteKeywords: ["BREAKING CHANGE", "BREAKING CHANGES", "BREAKING-CHANGE", "BREAKING"],

    // Incluir comparación con versión anterior
    issuePrefixes: ["#", "ISSUE-", "GH-"],

    // Scripts pre y post bump (opcional)
    scripts: {
        prebump: "go test ./...",
        postchangelog: "prettier --write CHANGELOG.md"
    }
};