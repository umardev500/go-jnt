if (window.transportInterval) {
    clearInterval(window.transportInterval);
}

document.title = "Penjadwalan Tugas Transportasi";

window.transportInterval = setInterval(() => {
    const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;

    document
        .querySelectorAll("table.el-table__body tr")
        .forEach(row => {

            const statusText = row.querySelectorAll("td")[10]?.innerText?.trim();
            const span = row.querySelectorAll("td")[23]?.querySelector("span");

            if (!span) return;

            // USED
            if (statusText === "Mengunggu proses" || statusText === "dalam perjalanan") {

                if (isDark) {
                    span.style.backgroundColor = "#92400e";
                    span.style.color = "#fef3c7";
                } else {
                    span.style.backgroundColor = "#fef3c7";
                    span.style.color = "#92400e";
                }

            }

            // UNUSED
            else {

                if (isDark) {
                    span.style.backgroundColor = "#1e3a8a";
                    span.style.color = "#f1f5f9";
                } else {
                    span.style.backgroundColor = "#e0f2fe";
                    span.style.color = "#0f172a";
                }
            }

            span.style.padding = "2px 6px";
            span.style.borderRadius = "6px";
            span.style.display = "inline-block";
        });

}, 1000);