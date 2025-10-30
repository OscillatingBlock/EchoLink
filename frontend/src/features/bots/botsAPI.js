const BASE_URL="https://api.echolink.studio/v1";

const authHeaders=() => {
    const token=localStorage.getItem('access_token');
    return token? {Authorization: `Bearer ${token}`} : {};
};

export const fetchBots=async () =>
{
    const res=await fetch(`${BASE_URL}/bots`, { headers: { ...authHeaders() } });
    if (!res.ok) throw new Error('Failed to fetch bots');
    return res.json();
};

export const createBot=async (payload) =>
{
    const res=await fetch(`${BASE_URL}/bots`, {
        method: "POST",
        headers: {"Content-Type": "application/json", ...authHeaders()},
        body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error('Failed to create bot');
    return res.json();
};

export const fetchBotDetails=async (id) =>
{
    const res=await fetch(`${BASE_URL}/bots/${id}`, { headers: { ...authHeaders() } });
    if (!res.ok) throw new Error('Failed to fetch bot');
    return res.json();
};

export const deleteBot=async (id) =>
{
    const res=await fetch(`${BASE_URL}/bots/${id}`, {
        method: "DELETE",
        headers: { ...authHeaders() },
    });
    if (!res.ok) throw new Error('Failed to delete bot');
    return res.json();
};
