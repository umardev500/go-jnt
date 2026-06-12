copy(
  Array.from(document.querySelectorAll("table.el-table__body tr"))
    .map(row => row.querySelector("td:nth-child(4)")?.innerText.trim())
    .filter(Boolean)
    .join("\n")
);


// Pelacakan
copy(
  Array.from(document.querySelectorAll("table.el-table__body tr"))
    .map(row => row.querySelector("td:nth-child(3)")?.innerText.trim())
    .filter(Boolean)
    .join("\n")
);


// Reg
copy(
  Array.from(document.querySelectorAll("table.el-table__body tr"))
    .map(row => row.querySelector("td:nth-child(3)")?.innerText.trim())
    .filter(text => text && !text.startsWith("JBGX"))
    .join("\n")
);

// with log
(() => {
  const result = Array.from(document.querySelectorAll("table.el-table__body tr"))
    .filter(row => {
      const status = row.querySelector("td:nth-child(7)")?.innerText.trim().toLowerCase();
      return status !== "dihapuskan";
    })
    .map(row => row.querySelector("td:nth-child(3)")?.innerText.trim())
    .filter(text => text && !text.startsWith("JBGX"));

  copy(result.join("\n"));

  console.log(`Copied ${result.length} rows`);
})();

// Only JBGX
copy(
  Array.from(document.querySelectorAll("table.el-table__body tr"))
    .map(row => row.querySelector("td:nth-child(3)")?.innerText.trim())
    .filter(text => text && text.startsWith("JBGX"))
    .join("\n")
);

// Exlude deleted
(() => {
  const result = Array.from(document.querySelectorAll("table.el-table__body tr"))
    .filter(row => {
      const status = row.querySelector("td:nth-child(7)")?.innerText.trim().toLowerCase();
      return status !== "dihapuskan";
    })
    .map(row => row.querySelector("td:nth-child(3)")?.innerText.trim())
    .filter(text => text && text.startsWith("JBGX"));

  copy(result.join("\n"));

  console.log(`Copied ${result.length} rows`);
})();

// 2026-05-04 21:59:00 nth-child 19
copy((() => {
  const results = Array.from(document.querySelectorAll("table.el-table__body tr"))
    .filter(row => {
      const tds = row.querySelectorAll("td");
      const timeText = tds[19]?.innerText.trim(); // ✅ correct column
      if (!timeText) return false;

      const time = timeText.split(" ")[1]; // "HH:MM:SS"
      return time >= "21:59:00" && time <= "22:00:00";
    })
    .map(row => {
      const tds = row.querySelectorAll("td");
      return tds[2]?.innerText.trim(); // 3rd column
    })
    .filter(Boolean);

  console.log("Copied count:", results.length);

  return results.join("\n");
})());

// skip status
copy((() => {
  const results = Array.from(document.querySelectorAll("table.el-table__body tr"))
    .filter(row => {
      const tds = row.querySelectorAll("td");

      const timeText = tds[19]?.innerText.trim();
      if (!timeText) return false;

      const status = tds[6]?.innerText.trim().toLowerCase();
      if (status === "dalam perjalanan" || status === "lengkap") return false;

      const time = timeText.split(" ")[1]; // "HH:MM:SS"
      return time >= "21:59:00" && time <= "22:00:00";
    })
    .map(row => {
      const tds = row.querySelectorAll("td");
      return tds[2]?.innerText.trim();
    })
    .filter(Boolean);

  console.log("Copied count:", results.length);

  return results.join("\n");
})());

// copy with detail
copy((() => {
  const results = Array.from(document.querySelectorAll("table.el-table__body tr"))
    .filter(row => {
      const tds = row.querySelectorAll("td");

      const timeText = tds[19]?.innerText.trim();
      if (!timeText) return false;

      const status = tds[6]?.innerText.trim().toLowerCase();
      if (status.includes("dalam perjalanan") || status.includes("lengkap")) return false;

      const time = timeText.split(" ")[1]; // "HH:MM:SS"
      return time >= "03:00:00" && time <= "08:01:00";
    })
    .map(row => {
      const tds = row.querySelectorAll("td");

      const col3 = tds[2]?.innerText.trim() || "";
      const col4 = tds[3]?.innerText.trim() || "";
      const date = tds[19]?.innerText.trim() || "";

      return `${col3} ${col4} ${date}`.trim();
    })
    .filter(Boolean);

  console.log("Copied count:", results.length);

  return results.join("\n");
})());