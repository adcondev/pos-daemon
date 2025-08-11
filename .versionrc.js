module.exports = {
    // Solo mostrar lo importante en el changelog
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "perf", section: "⚡ Performance"},
        {type: "deps", section: "📦 Dependencies"},
        {type: "revert", section: "⏪ Reverts"},
        // El resto se oculta
        {type: "docs", hidden: true},
        {type: "style", hidden: true},
        {type: "refactor", hidden: true},
        {type: "test", hidden: true},
        {type: "chore", hidden: true},
        {type: "ci", hidden: true},
        {type: "build", hidden: true}
    ],

    // Configuración de GitHub
    commitUrlFormat: "https://github.com/AdConDev/pos-daemon/commit/{{hash}}",
    compareUrlFormat: "https://github.com/AdConDev/pos-daemon/compare/{{previousTag}}...{{currentTag}}",

    // Skip CI en commits de release
    releaseCommitMessageFormat: "chore(release): v{{currentTag}} [skip ci]"
};