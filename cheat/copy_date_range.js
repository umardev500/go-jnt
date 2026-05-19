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