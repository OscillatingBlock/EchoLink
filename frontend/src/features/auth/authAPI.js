const BASE_URL = "https://api.echolink.studio/v1";

export const connectTwilio = async (payload) => {
  const res = await fetch(`${BASE_URL}/connect-twilio`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!res.ok) throw new Error("Failed to connect Twilio");
  return res.json();
};

export const getMyNumber = async () => {
  const token = localStorage.getItem("access_token");
  const res = await fetch(`${BASE_URL}/my-number`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error("Failed to fetch number");
  return res.json();
};
