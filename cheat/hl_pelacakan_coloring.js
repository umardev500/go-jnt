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

            // =========================
            // COLUMN VALUES
            // =========================
            const value3 = (
                tds[2]?.querySelector("span")?.textContent ||
                tds[2]?.textContent ||
                ""
            ).trim();

            const value6 = (
                tds[6]?.querySelector("span")?.textContent ||
                tds[6]?.textContent ||
                ""
            ).trim().toLowerCase();

            const value17 = (
                tds[17]?.querySelector("span")?.textContent ||
                tds[17]?.textContent ||
                ""
            ).trim();

            const span19 = tds[19]?.querySelector("span");

            // =========================
            // RESET ROW
            // =========================
            tds.forEach(td => {
                td.style.backgroundColor = "";
            });

            // =========================
            // PRIORITY 1 → TEMBAKAN
            // =========================
            if (value17 === "TEMBAKAN") {

                tds.forEach(td => {
                    td.style.backgroundColor = isDark
                        ? "#78350f"
                        : "#fff7ed";
                });

            }

            // =========================
            // PRIORITY 2 → JBGX
            // =========================
            else if (value3.startsWith("JBGX")) {

                tds.forEach(td => {
                    td.style.backgroundColor = isDark
                        ? "#14532d"
                        : "#dcfce7";
                });

            }

            // =========================
            // TD 19 STATUS COLOR
            // RULE BASED ON TD 6
            // =========================
            if (span19) {

                // ACTIVE / USED
                if (
                    value6 === "mengunggu proses" ||
                    value6 === "dalam perjalanan" ||
                    value6 === "lengkap"
                ) {

                    if (isDark) {
                        span19.style.backgroundColor = "#92400e";
                        span19.style.color = "#fef3c7";
                    } else {
                        span19.style.backgroundColor = "#fef3c7";
                        span19.style.color = "#92400e";
                    }

                } // DELETED
                else if (
                    value6 === "dihapuskan" ||
                    value6 === "deleted"
                ) {

                    if (isDark) {
                        span19.style.backgroundColor = "#7f1d1d";
                        span19.style.color = "#fecaca";
                    } else {
                        span19.style.backgroundColor = "#fee2e2";
                        span19.style.color = "#991b1b";
                    }

                }

                // DEFAULT / UNUSED
                else {

                    if (isDark) {
                        span19.style.backgroundColor = "#1e3a8a";
                        span19.style.color = "#f1f5f9";
                    } else {
                        span19.style.backgroundColor = "#e0f2fe";
                        span19.style.color = "#0f172a";
                    }
                }

                span19.style.padding = "2px 6px";
                span19.style.borderRadius = "6px";
                span19.style.display = "inline-block";
            }

        });

}, 1000);