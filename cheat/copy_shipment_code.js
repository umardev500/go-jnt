let result = document.querySelectorAll(".el-table__body-wrapper tbody tr");

let output = Array.from(result).map(tr => {
    let value = tr.children[2]?.innerText?.trim();
    return value ? "- " + value : null;
}).filter(Boolean);

console.log(output.join("\n"));