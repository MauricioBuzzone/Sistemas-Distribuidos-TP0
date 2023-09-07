import logging

class AgencyRegister:
    def __init__(self, num_agencies):
        # Initialize a list of booleans with the specified length, all set to False
        self.agency_status = [False] * num_agencies

    def update(self, agency_num):
        # Check if agency_num is within the valid range
        if 1 <= agency_num <= len(self.agency_status):
            # Update the agency's status to True
            self.agency_status[agency_num-1] = True
            return True
        else:
            return False

    def finish(self):
        # Return True if all booleans in the list are True
        return all(self.agency_status)