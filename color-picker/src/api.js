
export async function getStatus() {
  const resp = await fetch("/api/status/");
  return await resp.json();
}

export function setStatus(data) {
	return fetch(
		"/api/status/",
		{
			method: "POST",
			body: JSON.stringify(data),
			headers: {
					"Content-Type": "application/json",
			}
		},
	);
}
