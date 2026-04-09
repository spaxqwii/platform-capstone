import React, { useEffect, useState } from 'react';
import './App.css';

interface ServiceStatus {
  status: string;
  service?: string;
}

function App() {
  const [apiStatus, setApiStatus] = useState<string>('checking...');
  const [workerStatus, setWorkerStatus] = useState<string>('checking...');
  const [jobResult, setJobResult] = useState<string>('');

  const checkHealth = async () => {
    try {
      const apiRes = await fetch('http://localhost:8080/health');
      const apiData: ServiceStatus = await apiRes.json();
      setApiStatus(apiData.status);
    } catch {
      setApiStatus('down');
    }

    try {
      const workerRes = await fetch('http://localhost:8081/health');
      const workerData: ServiceStatus = await workerRes.json();
      setWorkerStatus(workerData.status);
    } catch {
      setWorkerStatus('down');
    }
  };

  const triggerJob = async () => {
    try {
      const res = await fetch('http://localhost:8081/job?data=test-job', { method: 'POST' });
      const data = await res.json();
      setJobResult(`Job queued: ${data.task_id}`);
    } catch (e) {
      setJobResult('Failed to queue job');
    }
  };

  useEffect(() => {
    checkHealth();
    const interval = setInterval(checkHealth, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Platform Dashboard</h1>
        
        <div style={{ margin: '20px' }}>
          <h2>Service Status</h2>
          <p>API: {apiStatus}</p>
          <p>Worker: {workerStatus}</p>
          <button onClick={checkHealth}>Refresh</button>
        </div>

        <div style={{ margin: '20px' }}>
          <h2>Job Control</h2>
          <button onClick={triggerJob}>Trigger Background Job</button>
          <p>{jobResult}</p>
        </div>
      </header>
    </div>
  );
}

export default App;