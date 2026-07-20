if (!window.toggleColumns) {
  window.toggleColumns = function(columns) {
    if (!columns.length) return;

    // Determine the action from the first column only
    const first = document.querySelector(
      `.el-table__header-wrapper th:nth-child(${columns[0]})`
    ) || document.querySelector(
      `.el-table__body-wrapper td:nth-child(${columns[0]})`
    );

    if (!first) return;

    const isHidden = getComputedStyle(first).display === "none";
    const newDisplay = isHidden ? "" : "none";

    columns.forEach(index => {
      document.querySelectorAll(
        `.el-table__header-wrapper th:nth-child(${index})`
      ).forEach(el => {
        el.style.display = newDisplay;
      });

      document.querySelectorAll(
        `.el-table__body-wrapper td:nth-child(${index})`
      ).forEach(el => {
        el.style.display = newDisplay;
      });
    });
  };
}

// Both columns will always hide/show together
// Penjadwalan
toggleColumns([6]);
toggleColumns([6, 7, 8, 9, 10, 12, 13, 14, 15, 16, 17, 19, 20, 21, 22, 23, 25, 28]);