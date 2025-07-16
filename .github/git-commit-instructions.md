Redacta el mensaje de commit siguiendo estrictamente el estándar “Conventional Commits 1.0.0”.

Formato:
<type>(<scope>): <short-summary>

<blank line>
<body – explicación en presente, 72 caracteres por línea máx>

<blank line>
<footer – refs externas, BREAKING CHANGE, Co-Authored-By, etc.>

Reglas:
1. Usa uno de estos <type>: feat | fix | docs | style | refactor | perf | test | build | ci | chore | revert.
2. <scope> debe ser breve (módulo, paquete o carpeta) y en kebab-case. Ej.: “printer-driver”.
3. <short-summary> en modo imperativo, ≤ 50 caracteres, sin punto final.
4. Si el cambio rompe compatibilidad añade la línea `BREAKING CHANGE:` en el footer.
5. El <body> debe responder “qué” y "por que", nunca “cómo” (el diff ya muestra el cómo) de forma resumida.
7. Multi-línea: envuelve a 72 caracteres; deja línea en blanco entre párrafos.
8. Mantén el mensaje en español salvo el <type> y palabras reservadas.
