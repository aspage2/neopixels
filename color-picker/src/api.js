
export async function getStatus() {
  const resp = await fetch("/api/lights/status");
  return await resp.json();
}

export function setStatus(data) {
	return fetch(
		"/api/lights/status",
		{
			method: "POST",
			body: JSON.stringify(data),
			headers: {
					"Content-Type": "application/json",
			}
		},
	);
}
