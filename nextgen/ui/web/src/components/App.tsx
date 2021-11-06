import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Header } from './header';
import { SettingsPage } from './settings';

export const App = () => {
  return (
    <BrowserRouter>
      <Header />
      <Routes>
        <Route path='/settings' element={<SettingsPage />}></Route>
      </Routes>
    </BrowserRouter>
  );
};
