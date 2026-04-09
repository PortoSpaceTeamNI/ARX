import { BrowserRouter, Route, Routes } from 'react-router-dom';
import 'dockview-react/dist/styles/dockview.css';

import FillingPage from '@/pages/FillingPage';
import { initWebSocket } from '@/utils/webSocketManager';

initWebSocket();

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/">
          <Route path="filling" element={<FillingPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
