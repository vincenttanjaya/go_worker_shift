const API_BASE = "http://localhost:8080";

export async function getAvailableShifts() {
  const res = await fetch(`${API_BASE}/shifts/available`);
  if (!res.ok) throw new Error('Failed to fetch');
  return await res.json();
}

export async function getRequests(workerId) {
  const res = await fetch(`${API_BASE}/worker/${workerId}/requests`);
  if (!res.ok) throw new Error('Failed to fetch');
  return await res.json();
}

export async function requestShift(workerId, shiftId) {
  const res = await fetch(`${API_BASE}/shift-requests`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ worker_id: workerId, shift_id: shiftId })
  });
  if (!res.ok) throw new Error('Failed to request shift');
}

export async function getWorkerById(workerId) {
  const res = await fetch(`${API_BASE}/workers/${workerId}`);
  if (!res.ok) throw new Error("Worker not found");
  return await res.json();
}