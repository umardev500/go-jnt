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