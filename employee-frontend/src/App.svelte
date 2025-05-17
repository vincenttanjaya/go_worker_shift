<script>
  import { getAvailableShifts, requestShift, getRequests, getWorkerById } from './api.js';

  let workerIdInput = '';
  let workerId = null;
  let worker = null;

  let shifts = [];
  let myRequests = [];
  let error = '';
  let loadingShifts = false;
  let loadingRequests = false;

  let dashboardReady = false;

  async function handleWorkerIdSubmit(e) {
    e.preventDefault();
    error = '';
    worker = null;
    dashboardReady = false;
    try {
      const idNum = parseInt(workerIdInput);
      if (!idNum || idNum < 1) { error = "Please enter a valid Worker ID"; return; }
      worker = await getWorkerById(idNum);
      workerId = idNum;
      dashboardReady = true;
      await loadData();
    } catch (e) {
      error = "Worker not found";
    }
  }

  async function loadData() {
    error = '';
    loadingShifts = loadingRequests = true;
    try {
      shifts = await getAvailableShifts();
    } catch (e) {
      error = "Could not fetch shifts";
      shifts = [];
    } finally {
      loadingShifts = false;
    }
    try {
      myRequests = await getRequests(workerId);
    } catch (e) {
      error = "Could not fetch requests";
      myRequests = [];
    } finally {
      loadingRequests = false;
    }
  }

  async function handleRequest(shiftId) {
    error = '';
    try {
      await requestShift(workerId, shiftId);
      await loadData();
    } catch (e) {
      error = "Failed to request shift (maybe already requested or assigned?)";
    }
  }

  // Check if worker already requested this shift
  function hasRequested(shiftId) {
    return myRequests.some(r => r.shift_id === shiftId);
  }
</script>

<style>
  body {
    font-family: sans-serif;
    background: #252525;
    color: #fafafa;
  }
  .error { color: #e74c3c; text-align: center; }
  .loading { color: #888; }
  h1, h2, h3 { text-align: center; }
  .centered { text-align: center; margin: 30px 0; }
  table { border-collapse: collapse; width: 90%; margin: 16px auto; }
  th, td { border: 1px solid #444; padding: 8px; }
  th { background: #222; color: #fff;}
  input[type="number"] { width: 100px; padding: 6px; font-size: 1rem; }
  button { padding: 5px 15px; font-size: 1rem; margin: 0 3px; }
</style>

<h1>Employee Shift Dashboard</h1>

<div class="centered">
  <form on:submit|preventDefault={handleWorkerIdSubmit}>
    <label>
      Enter Worker ID:
      <input type="number" min="1" bind:value={workerIdInput} required autofocus
        on:keydown={(e) => e.key === 'Enter' && handleWorkerIdSubmit(e)}/>
    </label>
    <button type="submit">Enter</button>
  </form>
  {#if error}
    <div class="error">{error}</div>
  {/if}
</div>

{#if dashboardReady && worker}
  <div class="centered">
    <h2>Hi, {worker.name} 👋</h2>
  </div>

  <h2 class="centered">Available Shifts</h2>
  {#if loadingShifts}
    <p class="loading centered">Loading shifts...</p>
  {:else if shifts.length === 0}
    <p class="centered">No available shifts</p>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Date</th>
          <th>Start</th>
          <th>End</th>
          <th>Role</th>
          <th>Location</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each shifts as s}
          <tr>
            <td>{s.date}</td>
            <td>{s.start_time}</td>
            <td>{s.end_time}</td>
            <td>{s.role}</td>
            <td>{s.location}</td>
            <td>
              {#if hasRequested(s.id)}
                <button disabled>Requested</button>
              {:else}
                <button on:click={() => handleRequest(s.id)}>Request</button>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}

  <h2 class="centered">My Shift Requests</h2>
  {#if loadingRequests}
    <p class="loading centered">Loading requests...</p>
  {:else if myRequests.length === 0}
    <p class="centered">You have not requested any shifts yet.</p>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Shift ID</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {#each myRequests as req}
          <tr>
            <td>{req.shift_id}</td>
            <td>{req.status}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
{/if}