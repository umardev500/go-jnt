(() => {
  const token = JSON.parse(localStorage.getItem("userData") || "{}")?.uuid;

  if (!token) {
    console.log("Token not found.");
    return;
  }

  copy(token);
  console.log("✅ Token copied.");
})();