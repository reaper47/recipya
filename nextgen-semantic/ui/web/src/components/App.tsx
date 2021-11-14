import { useState } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Container, Segment, Sidebar } from 'semantic-ui-react';
import { Header } from './header';
import { Drawer } from './header/Drawer';
import { SettingsPage } from './settings';

export const App = () => {
  const [isSidebarOpen, setIsSidebarVisible] = useState(false);

  const handleSidebar = (isOpen: boolean) => setIsSidebarVisible(isOpen);

  const sidebarCSS = {
    minHeight: '100vh',
    border: 'none',
    borderRadius: '0',
  };

  return (
    <BrowserRouter>
      <Sidebar.Pushable style={sidebarCSS}>
        <Header isSidebarOpen={isSidebarOpen} handleSidebar={handleSidebar} />
        <Drawer isSidebarOpen={isSidebarOpen} />

        <Sidebar.Pusher
          as={Container}
          style={{ maxWidth: 'calc(100vw - 10rem)' }}
        >
          <Routes>
            <Route path='/settings' element={<SettingsPage />}></Route>
          </Routes>
        </Sidebar.Pusher>
      </Sidebar.Pushable>
    </BrowserRouter>
  );
};
