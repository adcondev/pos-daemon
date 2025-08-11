module.exports = {
    // Solo mostrar lo importante en el changelog
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "perf", section: "⚡ Performance"},
        {type: "deps", section: "📦 Dependencies"},
        {type: "revert", section: "⏪ Reverts"},
        {type: "docs", section: "📝 Documentation"},
        {type: "style", section: "🎨 Styles"},
        {type: "refactor", section: "🔨 Refactoring"},
        {type: "test", section: "✅ Tests"},
        {type: "chore", hidden: true},
        {type: "ci", section: "🤖 Continuous Integration"},
        {type: "build", section: "🏗️ Build System"}
    ],

    // Configuración de GitHub
    commitUrlFormat: "https://github.com/AdConDev/pos-daemon/commit/{{hash}}",
    compareUrlFormat: "https://github.com/AdConDev/pos-daemon/compare/{{previousTag}}...{{currentTag}}",

    // Skip CI en commits de release
    releaseCommitMessageFormat: "chore(release): v{{currentTag}} [skip ci]"
};