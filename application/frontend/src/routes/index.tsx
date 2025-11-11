import { BrowserRouter, Routes, Route } from "react-router-dom";
import AppLayout from "../components/layout/AppLayout";

import Dashboard from "../pages/Dashboard";
import Agents from "../pages/Agents";
import Tokens from "../pages/Tokens";
import Offers from "../pages/Offers";
import Contracts from "../pages/Contracts";

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <AppLayout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/agents" element={<Agents />} />
          <Route path="/tokens" element={<Tokens />} />
          <Route path="/offers" element={<Offers />} />
          <Route path="/contracts" element={<Contracts />} />
        </Routes>
      </AppLayout>
    </BrowserRouter>
  );
}
