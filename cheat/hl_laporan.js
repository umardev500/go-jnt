document.title = "Laporan Pengiriman";

if (window.tableHighlightInterval) {
    clearInterval(window.tableHighlightInterval);
}

window.tableHighlightInterval = setInterval(() => {
    const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;

    document.querySelectorAll("table.el-table__body tr").forEach(row => {

        [34, 36, 38].forEach(colIndex => {
            const td = row.querySelectorAll("td")[colIndex];
            const innerSpan = td?.querySelector(".cell span span") || td?.querySelector("span span");

            if (!innerSpan) return;

            const text = innerSpan.textContent.trim();

            // outer wrapper (for styling)
            const wrapper = td.querySelector(".cell > span");
            if (!wrapper) return;

            // base style
            wrapper.style.padding = "2px 6px";
            wrapper.style.borderRadius = "6px";
            wrapper.style.display = "inline-block";

            if (isDark) {
                wrapper.style.backgroundColor = "#1e3a8a"; // dark blue
                wrapper.style.color = "#f1f5f9";           // light text
            } else {
                wrapper.style.backgroundColor = "#e0f2fe"; // light blue
                wrapper.style.color = "#0f172a";           // dark text
            }

            // RULE: detect ---
            if (text === "---" || text === "--") {

                if (colIndex === 36 || colIndex === 38) {
                    if (isDark) {
                        wrapper.style.backgroundColor = "#7f1d1d"; // dark red
                        wrapper.style.color = "#fee2e2";
                    } else {
                        wrapper.style.backgroundColor = "#fecaca"; // light red
                        wrapper.style.color = "#991b1b";
                    }
                } else {
                    if (isDark) {
                        wrapper.style.backgroundColor = "#374151"; // dark gray
                        wrapper.style.color = "#d1d5db";
                    } else {
                        wrapper.style.backgroundColor = "#e5e7eb"; // light gray
                        wrapper.style.color = "#6b7280";
                    }
                }
            }

        });

    });
}, 1000);

window.clickInterval = null;

function runClick() {
    // always clear first
    if (window.clickInterval) {
        clearInterval(window.clickInterval);
        window.clickInterval = null;
    }

    // run immediately
    const btn = document.querySelector(".btn-query");
    if (btn) btn.click();

    // then start interval again
    window.clickInterval = setInterval(() => {
        const btn = document.querySelector(".btn-query");
        if (btn) btn.click();
    }, 5000);
}

runClick();