<script>
  import { onMount } from 'svelte';
  import {
    getAllShifts, createShift, updateShift, deleteShift,
    getAllWorkers, createWorker, updateWorker, deleteWorker,
    getWorkerRequests, approveRequest, rejectRequest,
    getRequestsByStatus
  } from './api.js';

  let allShifts = [];
  let allWorkers = [];
  let pendingRequests = [];
  let error = '';

  let newShift = { date: "", start_time: "", end_time: "", role: "", location: "" };
  let editShift = null;

  function startEditShift(shift) {
    // Use a shallow copy and transform assigned_worker_id to a primitive for the select
    editShift = {
      ...shift,
      assigned_worker_id: (shift.assigned_worker_id && shift.assigned_worker_id.Valid)
        ? shift.assigned_worker_id.Int64
        : null
    };
  }
  function cancelEditShift() { editShift = null; }

  async function handleCreateShift() {
    error = '';
    try {
      await createShift(newShift);
      newShift = { date: "", start_time: "", end_time: "", role: "", location: "" };
      await loadData();
    } catch (e) { error = e.message; }
  }
  async function handleUpdateShift() {
    error = '';
    try {
      // Ensure we send assigned_worker_id as either null or a number
      let payload = {
        ...editShift,
        assigned_worker_id:
          (editShift.assigned_worker_id === "" || editShift.assigned_worker_id === null)
            ? null
            : Number(editShift.assigned_worker_id)
      };
      await updateShift(editShift.id, payload);
      editShift = null;
      await loadData();
    } catch (e) { error = e.message; }
  }
  async function handleDeleteShift(id) {
    if (window.confirm("Delete shift?")) {
      try { await deleteShift(id); await loadData(); }
      catch (e) { error = e.message; }
    }
  }

  // ---- WORKER ----
  let newWorkerName = '';
  let editWorker = null;

  function startEditWorker(w) { editWorker = { ...w }; }
  function cancelEditWorker() { editWorker = null; }

  async function handleCreateWorker() {
    error = '';
    try {
      await createWorker(newWorkerName);
      newWorkerName = "";
      await loadData();
    } catch (e) { error = e.message; }
  }
  async function handleUpdateWorker() {
    error = '';
    try {
      await updateWorker(editWorker.id, editWorker);
      editWorker = null;
      await loadData();
    } catch (e) { error = e.message; }
  }
  async function handleDeleteWorker(id) {
    if (window.confirm("Delete worker?")) {
      try { await deleteWorker(id); await loadData(); }
      catch (e) { error = e.message; }
    }
  }

  // ---- REQUESTS ----
  let allRequests = [];
  async function loadPendingRequests() {
    pendingRequests = await getRequestsByStatus("pending");
  }
  async function loadAllRequests() {
    let reqs = [];
    for (const w of allWorkers) {
      const wReqs = await getWorkerRequests(w.id);
      reqs.push(...wReqs);
    }
    allRequests = reqs;
  }

  async function handleApprove(req) {
    try { await approveRequest(req.id); await loadData(); }
    catch (e) { error = "Failed to approve"; }
  }
  async function handleReject(req) {
    try { await rejectRequest(req.id); await loadData(); }
    catch (e) { error = "Failed to reject"; }
  }

  async function loadData() {
    error = '';
    [allShifts, allWorkers] = await Promise.all([getAllShifts(), getAllWorkers()]);
    await loadPendingRequests();
    await loadAllRequests();
  }

  onMount(loadData);

  // ---------- Rendering helpers ----------
  // Used for displaying assigned worker in table
  function assignedWorkerName(obj) {
    if (!obj || !obj.Valid) return "Unassigned";
    const w = allWorkers.find(w => w.id === obj.Int64);
    return w ? w.name : `ID ${obj.Int64}`;
  }
</script>

<style>
  .error { color: red; }
  table { border-collapse: collapse; width: 95%; margin: 10px 0; }
  th, td { border: 1px solid #ccc; padding: 8px; }
  th { background: #222; color: #fff; }
  input, select { margin: 2px 4px; padding: 3px; }
  button { margin: 0 2px; }
  h1, h2, h3 { text-align: center; }
  .gray { color: #aaa; }
</style>

<h1>Admin Shift Management</h1>
{#if error}
  <div class="error">{error}</div>
{/if}

<!-- ===== Pending Shift Requests ===== -->
<h2>Pending Shift Requests</h2>
{#if pendingRequests.length === 0}
  <p style="text-align:center;">No pending shift requests.</p>
{:else}
  <table>
    <thead>
      <tr>
        <th>ID</th>
        <th>Worker</th>
        <th>Shift ID</th>
        <th>Status</th>
        <th>Action</th>
      </tr>
    </thead>
    <tbody>
      {#each pendingRequests as req}
        <tr>
          <td>{req.id}</td>
          <td>{#if allWorkers.length > 0}
            {#each allWorkers.filter(w => w.id === req.worker_id) as w}{w.name}{/each}
          {/if}</td>
          <td>{req.shift_id}</td>
          <td>{req.status}</td>
          <td>
            <button on:click={() => handleApprove(req)}>Approve</button>
            <button on:click={() => handleReject(req)}>Reject</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<!-- ===== All Shifts Table with Inline Add/Edit/Delete ===== -->
<h2>All Shifts</h2>
<table>
  <thead>
    <tr>
      <th>ID</th><th>Date</th><th>Start</th><th>End</th><th>Role</th><th>Location</th><th>Assigned Worker</th><th>Actions</th>
    </tr>
  </thead>
  <tbody>
    <!-- Add Shift Inline Row -->
    <tr>
      <td class="gray">new</td>
      <td><input type="date" bind:value={newShift.date} required /></td>
      <td><input type="time" bind:value={newShift.start_time} required /></td>
      <td><input type="time" bind:value={newShift.end_time} required /></td>
      <td><input bind:value={newShift.role} required /></td>
      <td><input bind:value={newShift.location} /></td>
      <td class="gray">-</td>
      <td>
        <button on:click={handleCreateShift}
          disabled={!newShift.date || !newShift.start_time || !newShift.end_time || !newShift.role}>Add</button>
      </td>
    </tr>
    {#each allShifts as s (s.id)}
      {#if editShift && editShift.id === s.id}
        <tr>
          <td>{editShift.id}</td>
          <td><input type="date" bind:value={editShift.date} required /></td>
          <td><input type="time" bind:value={editShift.start_time} required /></td>
          <td><input type="time" bind:value={editShift.end_time} required /></td>
          <td><input bind:value={editShift.role} required /></td>
          <td><input bind:value={editShift.location} /></td>
          <td>
            <select
              on:change={(e) => {
                const val = e.target.value;
                editShift.assigned_worker_id = val === "" ? null : Number(val);
              }}
            >
              <option value="">Unassigned</option>
              {#each allWorkers as w}
                <option value={w.id} selected={editShift.assigned_worker_id == w.id}>{w.name}</option>
              {/each}
            </select>
          </td>
          <td>
            <button on:click={handleUpdateShift}>Save</button>
            <button on:click={cancelEditShift}>Cancel</button>
          </td>
        </tr>
      {:else}
        <tr>
          <td>{s.id}</td>
          <td>{s.date}</td>
          <td>{s.start_time}</td>
          <td>{s.end_time}</td>
          <td>{s.role}</td>
          <td>{s.location}</td>
          <td>{assignedWorkerName(s.assigned_worker_id)}</td>
          <td>
            <button on:click={() => startEditShift(s)}>Edit</button>
            <button on:click={() => handleDeleteShift(s.id)}>Delete</button>
          </td>
        </tr>
      {/if}
    {/each}
  </tbody>
</table>

<!-- ===== All Workers Table with Inline Add/Edit/Delete ===== -->
<h2>All Workers</h2>
<table>
  <thead>
    <tr><th>ID</th><th>Name</th><th>Actions</th></tr>
  </thead>
  <tbody>
    <!-- Add Worker Inline Row -->
    <tr>
      <td class="gray">new</td>
      <td>
        <input bind:value={newWorkerName} required />
      </td>
      <td>
        <button on:click={handleCreateWorker} disabled={!newWorkerName}>Add</button>
      </td>
    </tr>
    {#each allWorkers as w (w.id)}
      {#if editWorker && editWorker.id === w.id}
        <tr>
          <td>{w.id}</td>
          <td><input bind:value={editWorker.name} /></td>
          <td>
            <button on:click={handleUpdateWorker}>Save</button>
            <button on:click={cancelEditWorker}>Cancel</button>
          </td>
        </tr>
      {:else}
        <tr>
          <td>{w.id}</td>
          <td>{w.name}</td>
          <td>
            <button on:click={() => startEditWorker(w)}>Edit</button>
            <button on:click={() => handleDeleteWorker(w.id)}>Delete</button>
          </td>
        </tr>
      {/if}
    {/each}
  </tbody>
</table>