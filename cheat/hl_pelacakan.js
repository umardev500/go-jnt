document.title = "Pelacakan Pengiriman";

if (window.tableInterval) {
    clearInterval(window.tableInterval);
}

window.tableInterval = setInterval(() => {
    const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;

    document
        .querySelectorAll("table.el-table__body tr")
        .forEach(row => {

            const tds = row.querySelectorAll("td");

            // ===== PREFIX RULE (td 3 → index 2) =====
            const cell3 = tds[2]?.querySelector("span") || tds[2];
            const value3 = cell3?.textContent.trim();

            if (value3 && value3.startsWith("JBGX")) {
                tds.forEach(td => {
                    td.style.backgroundColor = isDark
                        ? "#14532d"   // dark green
                        : "#dcfce7";  // light green
                });
            }

            // ===== ROW RULE (td 17) =====
            const cell17 = tds[17]?.querySelector("span");
            const value17 = cell17?.textContent.trim();

            if (value17 === "TEMBAKAN") {
                tds.forEach(td => {
                    td.style.backgroundColor = isDark
                        ? "#78350f"
                        : "#fff7ed";
                });
            } else if (!value3?.startsWith("JBGX")) {
                // reset only if NOT matched by JBGX rule
                tds.forEach(td => {
                    td.style.backgroundColor = "";
                });
            }

            // ===== CELL RULE (td 19) =====
            const span = tds[19]?.querySelector("span");

            if (span && !span.dataset.styled) {
                if (isDark) {
                    span.style.backgroundColor = "#1e3a8a";
                    span.style.color = "#f1f5f9";
                } else {
                    span.style.backgroundColor = "#e0f2fe";
                    span.style.color = "#0f172a";
                }

                span.style.padding = "2px 6px";
                span.style.borderRadius = "6px";
                span.style.display = "inline-block";

                span.dataset.styled = "true";
            }
        });
}, 1000);