module.exports = {
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "docs", section: "📚 Documentation"},
        {type: "chore", section: "🔧 Maintenance"},
        {type: "style", section: "💄 Style Changes"},
        {type: "refactor", section: "♻️ Code Refactoring"},
        {type: "perf", section: "⚡ Performance Improvements"},
        {type: "test", section: "✅ Testing"}
    ],
    commitUrlFormat: "https://github.com/AdConDev/yalemi-api/commit/{{hash}}",
    compareUrlFormat: "https://github.com/AdConDev/yalemi-api/compare/{{previousTag}}...{{currentTag}}",
    releaseCommitMessageFormat: "chore(release): {{currentTag}} [skip ci]",
    scopes: [
        "encoding",
        "connector",
        "printer",
        "imaging",
        "profiles",
        "protocol",
        "escpos",
        "errors",
        "logs",
        "config",
        "api",
        "models",
        "service",
        "utils"
    ],
    noteKeywords: [
        "BREAKING CHANGE",
        "BREAKING CHANGES",
        "BREAKING"
    ]
};