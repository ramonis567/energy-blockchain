import unittest
from unittest.mock import patch, MagicMock
import sys
import os

# Add the src directory to the Python path
src_path = os.path.abspath(os.path.join(os.path.dirname(__file__), 'src'))
sys.path.insert(0, src_path)

# Mock modules that cannot be installed
sys.modules['paho'] = MagicMock()
sys.modules['paho.mqtt'] = MagicMock()
sys.modules['paho.mqtt.client'] = MagicMock()

# Mock config to avoid import errors
sys.modules['config'] = MagicMock()

from simulation import simulate

class TestSimulation(unittest.TestCase):

    @patch('simulation.publish_data')
    @patch('simulation.connect')
    @patch('simulation.disconnect')
    @patch('simulation.time.sleep')
    @patch('builtins.open', new_callable=unittest.mock.mock_open, read_data='{"users":[{"id": "1", "type": "consumer", "class": "residential"}]}')
    def test_simulate_runs_once(self, mock_open, mock_sleep, mock_disconnect, mock_connect, mock_publish):
        # We want the simulation to run only once for this test
        mock_sleep.side_effect = InterruptedError

        try:
            simulate()
        except InterruptedError:
            pass # Expected interruption

        # Check that publish_data was called
        self.assertTrue(mock_publish.called)

        # Check that connect was called
        self.assertTrue(mock_connect.called)

        # Check that disconnect was called (due to InterruptedError handled in finally block)
        self.assertTrue(mock_disconnect.called)

if __name__ == '__main__':
    unittest.main()
