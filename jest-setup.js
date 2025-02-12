// Jest setup provided by Grafana scaffolding
import './.config/jest-setup';

// mock the intersection observer and just say everything is in view
const mockIntersectionObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));
global.IntersectionObserver = mockIntersectionObserver;
