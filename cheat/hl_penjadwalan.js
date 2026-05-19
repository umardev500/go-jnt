if (window.transportInterval) {
    clearInterval(window.transportInterval);
}

document.title = "Penjadwalan Tugas Transportasi";

// Start fresh interval
window.transportInterval = setInterval(() => {
    const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;

    document
        .querySelectorAll("table.el-table__body tr")
        .forEach(row => {
            const span = row.querySelectorAll("td")[23]?.querySelector("span");

            if (span && !span.dataset.styled) {
                if (isDark) {
                    span.style.backgroundColor = "#1e3a8a"; // dark blue
                    span.style.color = "#f1f5f9";           // light text
                } else {
                    span.style.backgroundColor = "#e0f2fe"; // light blue
                    span.style.color = "#0f172a";           // dark text
                }

                span.style.padding = "2px 6px";
                span.style.borderRadius = "6px";
                span.style.display = "inline-block";

                span.dataset.styled = "true";
            }
        });
}, 1000);