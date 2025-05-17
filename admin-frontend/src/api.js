const API_BASE = "http://localhost:8080";

export async function getAllShifts() {
  const res = await fetch(`${API_BASE}/shifts`);
  return await res.json();
}

export async function getAssignedShifts() {
  const res = await fetch(`${API_BASE}/shifts/assigned`);
  return await res.json();
}

export async function getAvailableShifts() {
  const res = await fetch(`${API_BASE}/shifts/available`);
  return await res.json();
}

export async function getWorkerRequests(workerId) {
  const res = await fetch(`${API_BASE}/worker/${workerId}/requests`);
  return await res.json();
}

export async function approveRequest(requestId) {
  const res = await fetch(`${API_BASE}/shift-requests/${requestId}/approve`, { method: 'POST' });
  if (!res.ok) throw new Error('Unable to approve request');
}

export async function rejectRequest(requestId) {
  const res = await fetch(`${API_BASE}/shift-requests/${requestId}/reject`, { method: 'POST' });
  if (!res.ok) throw new Error('Unable to reject request');
}

export async function getAllWorkers() {
  const res = await fetch(`${API_BASE}/workers`);
  return await res.json();
}

export async function createShift(shift) {
    const res = await fetch(`${API_BASE}/shifts`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(shift)
    });
    if (!res.ok) throw new Error("Failed to create shift");
}

export async function updateShift(shiftId, updateFields) {
    const res = await fetch(`http://localhost:8080/shifts/${shiftId}`, {
      method: 'PUT',
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(updateFields),
    });
    if (!res.ok) throw new Error('Failed to update shift');
}
  
export async function deleteShift(shiftId) {
const res = await fetch(`http://localhost:8080/shifts/${shiftId}`, { method: 'DELETE' });
if (!res.ok) throw new Error('Failed to delete shift');
}

  export async function createWorker(name) {
    const res = await fetch(`${API_BASE}/signup`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    });
    if (!res.ok) throw new Error("Failed to create worker");
  }
  export async function updateWorker(id, worker) {
    const res = await fetch(`${API_BASE}/workers/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(worker)
    });
    if (!res.ok) throw new Error("Failed to update worker");
  }
  export async function deleteWorker(id) {
    const res = await fetch(`${API_BASE}/workers/${id}`, {
      method: "DELETE"
    });
    if (!res.ok) throw new Error("Failed to delete worker");
  }
  
export async function getRequestsByStatus(status) {
  const url = `${API_BASE}/shift-requests?status=${status}`;
  const res = await fetch(url, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' }
  });
  return await res.json();
}