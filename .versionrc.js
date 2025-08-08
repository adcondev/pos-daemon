module.exports = {
    // Configuración minimalista pero efectiva
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "perf", section: "⚡ Performance"},
        {type: "deps", section: "📦 Dependencies"},
        // Ocultar tipos que generan ruido en el changelog
        {type: "docs", section: "📚 Documentation", hidden: true},
        {type: "style", hidden: true},
        {type: "refactor", hidden: true},
        {type: "test", hidden: true},
        {type: "chore", hidden: true},
        {type: "ci", hidden: true},
        {type: "build", hidden: true},
        {type: "revert", section: "⏪ Reverts"}
    ],

    // URLs automáticas para GitHub
    commitUrlFormat: "https://github.com/{{owner}}/{{repository}}/commit/{{hash}}",
    compareUrlFormat: "https://github.com/{{owner}}/{{repository}}/compare/{{previousTag}}...{{currentTag}}",
    issueUrlFormat: "https://github.com/{{owner}}/{{repository}}/issues/{{id}}",

    // Configuración simple
    releaseCommitMessageFormat: "chore(release): v{{currentTag}} [skip ci]",
    skip: {
        commit: false,
        tag: false
    }
};