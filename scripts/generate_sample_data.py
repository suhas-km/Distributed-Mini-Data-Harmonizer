#!/usr/bin/env python3
"""
Mock Healthcare Data Generator

This script generates realistic mock healthcare data for testing and development.
It creates multiple CSV files representing different aspects of healthcare data:
- patients.csv: Patient demographic information
- vitals.csv: Patient vital signs and measurements
- medications.csv: Medication prescriptions and history
- lab_results.csv: Laboratory test results

The data includes intentional inconsistencies, duplicates, and format variations
to simulate real-world data harmonization challenges.
"""

import os
import csv
import random
import uuid
import datetime
import argparse
from typing import List, Dict, Any, Tuple

# Constants for generating realistic data
GENDERS = ["Male", "Female", "Non-binary", "Other", "Unknown", "M", "F"]
STATES = ["AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA", 
          "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", 
          "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", 
          "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", 
          "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY"]
BLOOD_TYPES = ["A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"]
ETHNICITIES = ["Hispanic or Latino", "Not Hispanic or Latino", "Unknown", "Declined"]
RACES = ["White", "Black or African American", "Asian", "American Indian or Alaska Native",
         "Native Hawaiian or Other Pacific Islander", "Other", "Multiple", "Unknown", "Declined"]
MARITAL_STATUS = ["Single", "Married", "Divorced", "Widowed", "Separated", "Unknown"]
INSURANCE_TYPES = ["Medicare", "Medicaid", "Private", "Self-pay", "Other"]

# Vital signs
VITAL_TYPES = ["Blood Pressure", "Heart Rate", "Respiratory Rate", "Temperature", "Oxygen Saturation", "Weight", "Height"]
VITAL_UNITS = {
    "Blood Pressure": "mmHg",
    "Heart Rate": "bpm",
    "Respiratory Rate": "breaths/min",
    "Temperature": "°C",
    "Oxygen Saturation": "%",
    "Weight": "kg",
    "Height": "cm"
}
VITAL_RANGES = {
    "Blood Pressure": lambda: f"{random.randint(90, 180)}/{random.randint(60, 110)}",
    "Heart Rate": lambda: str(random.randint(40, 120)),
    "Respiratory Rate": lambda: str(random.randint(12, 25)),
    "Temperature": lambda: f"{round(random.uniform(36.0, 39.0), 1)}",
    "Oxygen Saturation": lambda: str(random.randint(88, 100)),
    "Weight": lambda: f"{round(random.uniform(45.0, 120.0), 1)}",
    "Height": lambda: str(random.randint(150, 200))
}

# Medications
MEDICATIONS = [
    "Lisinopril", "Atorvastatin", "Levothyroxine", "Metformin", "Amlodipine",
    "Metoprolol", "Omeprazole", "Simvastatin", "Losartan", "Albuterol",
    "Gabapentin", "Hydrochlorothiazide", "Sertraline", "Montelukast", "Pantoprazole"
]
DOSAGES = ["5mg", "10mg", "20mg", "25mg", "50mg", "75mg", "100mg", "150mg", "200mg", "250mg", "500mg"]
FREQUENCIES = ["Once daily", "Twice daily", "Three times daily", "Four times daily", 
               "Every morning", "Every evening", "As needed", "Weekly", "Monthly"]

# Lab tests
LAB_TESTS = [
    "Complete Blood Count", "Basic Metabolic Panel", "Comprehensive Metabolic Panel",
    "Lipid Panel", "Thyroid Stimulating Hormone", "Hemoglobin A1C", "Urinalysis",
    "Liver Function Tests", "Blood Glucose", "Prothrombin Time"
]
LAB_UNITS = {
    "Complete Blood Count": "cells/μL",
    "Basic Metabolic Panel": "mmol/L",
    "Comprehensive Metabolic Panel": "mg/dL",
    "Lipid Panel": "mg/dL",
    "Thyroid Stimulating Hormone": "mIU/L",
    "Hemoglobin A1C": "%",
    "Urinalysis": "",
    "Liver Function Tests": "U/L",
    "Blood Glucose": "mg/dL",
    "Prothrombin Time": "sec"
}
ABNORMAL_FLAGS = ["Normal", "High", "Low", "Critical High", "Critical Low", "N", "H", "L", ""]

# Provider information
PROVIDER_PREFIXES = ["Dr.", ""]
PROVIDER_SPECIALTIES = ["Family Medicine", "Internal Medicine", "Cardiology", "Endocrinology", 
                       "Gastroenterology", "Neurology", "Oncology", "Pediatrics", "Psychiatry"]

def generate_patient_id() -> str:
    """Generate a patient ID with different formats to simulate data inconsistency"""
    format_type = random.randint(1, 4)
    if format_type == 1:
        return f"P{random.randint(10000, 99999)}"
    elif format_type == 2:
        return str(uuid.uuid4())[:8]
    elif format_type == 3:
        return f"PT-{random.randint(1000, 9999)}-{random.randint(10, 99)}"
    else:
        return f"{random.randint(100000, 999999)}"

def generate_name() -> Tuple[str, str, str]:
    """Generate a random name"""
    first_names = ["James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda", 
                  "William", "Elizabeth", "David", "Susan", "Richard", "Jessica", "Joseph", "Sarah",
                  "Thomas", "Karen", "Charles", "Nancy", "Christopher", "Lisa", "Daniel", "Margaret",
                  "Matthew", "Betty", "Anthony", "Sandra", "Mark", "Ashley", "Donald", "Dorothy",
                  "Steven", "Kimberly", "Paul", "Emily", "Andrew", "Donna", "Joshua", "Michelle",
                  "Kenneth", "Carol", "Kevin", "Amanda", "Brian", "Melissa", "George", "Deborah"]
    
    last_names = ["Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson",
                 "Moore", "Taylor", "Anderson", "Thomas", "Jackson", "White", "Harris", "Martin",
                 "Thompson", "Garcia", "Martinez", "Robinson", "Clark", "Rodriguez", "Lewis", "Lee",
                 "Walker", "Hall", "Allen", "Young", "Hernandez", "King", "Wright", "Lopez",
                 "Hill", "Scott", "Green", "Adams", "Baker", "Gonzalez", "Nelson", "Carter",
                 "Mitchell", "Perez", "Roberts", "Turner", "Phillips", "Campbell", "Parker", "Evans"]
    
    middle_initials = ["A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", 
                      "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", ""]
    
    first_name = random.choice(first_names)
    middle_name = random.choice(middle_initials)
    last_name = random.choice(last_names)
    
    return first_name, middle_name, last_name

def generate_address() -> Tuple[str, str, str, str]:
    """Generate a random address"""
    street_numbers = [str(random.randint(1, 9999)) for _ in range(100)]
    street_names = ["Main", "Oak", "Pine", "Maple", "Cedar", "Elm", "Washington", "Lake", 
                   "Hill", "Park", "Spring", "River", "Meadow", "Forest", "Highland", "Valley",
                   "Madison", "Jefferson", "Adams", "Monroe", "Lincoln", "Franklin", "Clinton"]
    street_types = ["St", "Ave", "Blvd", "Dr", "Ln", "Rd", "Way", "Pl", "Ct", "Terrace"]
    cities = ["Springfield", "Franklin", "Greenville", "Bristol", "Clinton", "Salem", "Madison",
             "Georgetown", "Arlington", "Fairview", "Riverside", "Centerville", "Manchester",
             "Auburn", "Dayton", "Lexington", "Oxford", "Burlington", "Milton", "Newport"]
    
    street = f"{random.choice(street_numbers)} {random.choice(street_names)} {random.choice(street_types)}"
    city = random.choice(cities)
    state = random.choice(STATES)
    zip_code = f"{random.randint(10000, 99999)}"
    
    return street, city, state, zip_code

def generate_phone() -> str:
    """Generate a phone number with different formats to simulate inconsistency"""
    format_type = random.randint(1, 4)
    area_code = random.randint(200, 999)
    prefix = random.randint(200, 999)
    line = random.randint(1000, 9999)
    
    if format_type == 1:
        return f"({area_code}) {prefix}-{line}"
    elif format_type == 2:
        return f"{area_code}-{prefix}-{line}"
    elif format_type == 3:
        return f"{area_code}.{prefix}.{line}"
    else:
        return f"{area_code}{prefix}{line}"

def generate_email(first_name: str, last_name: str) -> str:
    """Generate an email address based on name"""
    domains = ["gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "aol.com", "icloud.com"]
    
    # Sometimes include middle initial or numbers
    if random.random() < 0.3:
        username = f"{first_name.lower()}.{last_name.lower()}"
    elif random.random() < 0.6:
        username = f"{first_name.lower()}{last_name.lower()}{random.randint(1, 99)}"
    else:
        username = f"{first_name[0].lower()}{last_name.lower()}"
    
    # Sometimes don't include email
    if random.random() < 0.1:
        return ""
    
    return f"{username}@{random.choice(domains)}"

def generate_date_of_birth() -> str:
    """Generate a random date of birth for an adult patient"""
    year = random.randint(1940, 2005)
    month = random.randint(1, 12)
    day = random.randint(1, 28)  # Simplified to avoid month length issues
    
    # Format with different date formats to simulate inconsistency
    format_type = random.randint(1, 3)
    if format_type == 1:
        return f"{month:02d}/{day:02d}/{year}"
    elif format_type == 2:
        return f"{year}-{month:02d}-{day:02d}"
    else:
        return f"{day:02d}-{month:02d}-{year}"

def generate_provider() -> Tuple[str, str]:
    """Generate a provider name and ID"""
    first, _, last = generate_name()
    prefix = random.choice(PROVIDER_PREFIXES)
    provider_name = f"{prefix} {first} {last}"
    provider_id = f"PROV{random.randint(10000, 99999)}"
    return provider_name, provider_id

def generate_datetime(start_date: datetime.datetime, end_date: datetime.datetime) -> str:
    """Generate a random datetime between start and end dates"""
    delta = end_date - start_date
    random_days = random.randint(0, delta.days)
    random_date = start_date + datetime.timedelta(days=random_days)
    
    # Add random hours and minutes
    random_date = random_date.replace(
        hour=random.randint(8, 17),
        minute=random.choice([0, 15, 30, 45])
    )
    
    # Format with different datetime formats to simulate inconsistency
    format_type = random.randint(1, 3)
    if format_type == 1:
        return random_date.strftime("%Y-%m-%d %H:%M:%S")
    elif format_type == 2:
        return random_date.strftime("%m/%d/%Y %I:%M %p")
    else:
        return random_date.strftime("%d-%b-%Y %H:%M")

def generate_patient_data(num_patients: int) -> List[Dict[str, Any]]:
    """Generate a list of patient data dictionaries"""
    patients = []
    
    # Create a set to track used patient IDs to avoid duplicates
    # But intentionally create some duplicates to simulate data issues
    used_patient_ids = set()
    
    for _ in range(num_patients):
        first_name, middle_name, last_name = generate_name()
        street, city, state, zip_code = generate_address()
        
        # Generate patient ID, occasionally reusing an existing one
        if random.random() < 0.05 and used_patient_ids:  # 5% chance of duplicate
            patient_id = random.choice(list(used_patient_ids))
        else:
            patient_id = generate_patient_id()
            used_patient_ids.add(patient_id)
        
        dob = generate_date_of_birth()
        
        # Sometimes use different formats for gender
        gender_format = random.randint(1, 3)
        if gender_format == 1:
            gender = random.choice(GENDERS)
        elif gender_format == 2:
            gender = random.choice(["male", "female", "non-binary", "other", "unknown"])
        else:
            gender = random.choice(["M", "F", "NB", "O", "U"])
        
        provider_name, provider_id = generate_provider()
        
        patient = {
            "patient_id": patient_id,
            "first_name": first_name,
            "middle_name": middle_name,
            "last_name": last_name,
            "date_of_birth": dob,
            "gender": gender,
            "address": street,
            "city": city,
            "state": state,
            "zip_code": zip_code,
            "phone_number": generate_phone(),
            "email": generate_email(first_name, last_name),
            "blood_type": random.choice(BLOOD_TYPES) if random.random() > 0.2 else "",
            "ethnicity": random.choice(ETHNICITIES),
            "race": random.choice(RACES),
            "marital_status": random.choice(MARITAL_STATUS),
            "insurance_type": random.choice(INSURANCE_TYPES),
            "insurance_id": f"INS-{random.randint(100000, 999999)}",
            "primary_care_provider": provider_name,
            "provider_id": provider_id
        }
        
        # Introduce some missing data
        if random.random() < 0.1:  # 10% chance of missing email
            patient["email"] = ""
        if random.random() < 0.05:  # 5% chance of missing phone
            patient["phone_number"] = ""
        if random.random() < 0.15:  # 15% chance of missing blood type
            patient["blood_type"] = ""
        
        patients.append(patient)
    
    return patients

def generate_vitals_data(patients: List[Dict[str, Any]], num_records: int) -> List[Dict[str, Any]]:
    """Generate vital signs data for patients"""
    vitals = []
    start_date = datetime.datetime.now() - datetime.timedelta(days=365)  # Last year
    end_date = datetime.datetime.now()
    
    for _ in range(num_records):
        patient = random.choice(patients)
        vital_type = random.choice(VITAL_TYPES)
        
        vital = {
            "patient_id": patient["patient_id"],
            "observation_datetime": generate_datetime(start_date, end_date),
            "vital_type": vital_type,
            "vital_value": VITAL_RANGES[vital_type](),
            "unit": VITAL_UNITS[vital_type],
            "provider_id": patient["provider_id"] if random.random() > 0.3 else generate_provider()[1]
        }
        
        # Sometimes use different formats for vital types
        if random.random() < 0.1:
            if vital_type == "Blood Pressure":
                vital["vital_type"] = "BP"
            elif vital_type == "Heart Rate":
                vital["vital_type"] = "HR"
            elif vital_type == "Respiratory Rate":
                vital["vital_type"] = "RR"
            elif vital_type == "Temperature":
                vital["vital_type"] = "Temp"
            elif vital_type == "Oxygen Saturation":
                vital["vital_type"] = "O2 Sat"
        
        vitals.append(vital)
    
    return vitals

def generate_medication_data(patients: List[Dict[str, Any]], num_records: int) -> List[Dict[str, Any]]:
    """Generate medication data for patients"""
    medications = []
    start_date = datetime.datetime.now() - datetime.timedelta(days=365)  # Last year
    end_date = datetime.datetime.now() + datetime.timedelta(days=180)  # 6 months into future for end dates
    
    for _ in range(num_records):
        patient = random.choice(patients)
        medication_name = random.choice(MEDICATIONS)
        start_datetime = generate_datetime(start_date, datetime.datetime.now())
        
        # End date might be in the future or empty (ongoing)
        has_end_date = random.random() > 0.4  # 60% chance of having an end date
        if has_end_date:
            # Parse the start date to ensure end date is after
            try:
                if "-" in start_datetime:
                    if " " in start_datetime:
                        if "%p" in start_datetime:
                            start_dt = datetime.datetime.strptime(start_datetime, "%Y-%m-%d %I:%M %p")
                        else:
                            start_dt = datetime.datetime.strptime(start_datetime, "%Y-%m-%d %H:%M:%S")
                    else:
                        start_dt = datetime.datetime.strptime(start_datetime, "%Y-%m-%d")
                elif "/" in start_datetime:
                    if " " in start_datetime:
                        start_dt = datetime.datetime.strptime(start_datetime, "%m/%d/%Y %I:%M %p")
                    else:
                        start_dt = datetime.datetime.strptime(start_datetime, "%m/%d/%Y")
                else:
                    if " " in start_datetime:
                        start_dt = datetime.datetime.strptime(start_datetime, "%d-%b-%Y %H:%M")
                    else:
                        start_dt = datetime.datetime.strptime(start_datetime, "%d-%m-%Y")
            except ValueError:
                # If parsing fails, use current date as fallback
                start_dt = datetime.datetime.now()
            
            # Generate end date after start date
            end_datetime = generate_datetime(start_dt, end_date)
        else:
            end_datetime = ""
        
        medication = {
            "patient_id": patient["patient_id"],
            "medication_name": medication_name,
            "dosage": random.choice(DOSAGES),
            "frequency": random.choice(FREQUENCIES),
            "start_date": start_datetime,
            "end_date": end_datetime,
            "prescriber_id": patient["provider_id"] if random.random() > 0.3 else generate_provider()[1]
        }
        
        # Introduce some data inconsistencies
        if random.random() < 0.1:  # 10% chance of different medication name format
            medication["medication_name"] = medication_name.upper()
        if random.random() < 0.1:  # 10% chance of abbreviated frequency
            if medication["frequency"] == "Once daily":
                medication["frequency"] = "QD"
            elif medication["frequency"] == "Twice daily":
                medication["frequency"] = "BID"
            elif medication["frequency"] == "Three times daily":
                medication["frequency"] = "TID"
            elif medication["frequency"] == "Four times daily":
                medication["frequency"] = "QID"
        
        medications.append(medication)
    
    return medications

def generate_lab_data(patients: List[Dict[str, Any]], num_records: int) -> List[Dict[str, Any]]:
    """Generate lab result data for patients"""
    lab_results = []
    start_date = datetime.datetime.now() - datetime.timedelta(days=365)  # Last year
    end_date = datetime.datetime.now()
    
    for _ in range(num_records):
        patient = random.choice(patients)
        test_name = random.choice(LAB_TESTS)
        
        # Generate result value based on test type
        if test_name == "Complete Blood Count":
            result_value = str(round(random.uniform(3.5, 6.0), 1))
        elif test_name == "Basic Metabolic Panel":
            result_value = str(round(random.uniform(135, 145), 1))
        elif test_name == "Comprehensive Metabolic Panel":
            result_value = str(round(random.uniform(70, 110), 1))
        elif test_name == "Lipid Panel":
            result_value = str(random.randint(120, 240))
        elif test_name == "Thyroid Stimulating Hormone":
            result_value = str(round(random.uniform(0.4, 4.0), 2))
        elif test_name == "Hemoglobin A1C":
            result_value = str(round(random.uniform(4.0, 9.0), 1))
        elif test_name == "Urinalysis":
            result_value = random.choice(["Negative", "Positive", "Trace", "+1", "+2", "+3"])
        elif test_name == "Liver Function Tests":
            result_value = str(random.randint(10, 80))
        elif test_name == "Blood Glucose":
            result_value = str(random.randint(70, 180))
        elif test_name == "Prothrombin Time":
            result_value = str(round(random.uniform(10.0, 14.0), 1))
        
        # Generate reference range
        if test_name == "Complete Blood Count":
            ref_range = "3.5-5.5"
        elif test_name == "Basic Metabolic Panel":
            ref_range = "135-145"
        elif test_name == "Comprehensive Metabolic Panel":
            ref_range = "70-99"
        elif test_name == "Lipid Panel":
            ref_range = "<200"
        elif test_name == "Thyroid Stimulating Hormone":
            ref_range = "0.4-4.0"
        elif test_name == "Hemoglobin A1C":
            ref_range = "4.0-5.6"
        elif test_name == "Urinalysis":
            ref_range = "Negative"
        elif test_name == "Liver Function Tests":
            ref_range = "10-40"
        elif test_name == "Blood Glucose":
            ref_range = "70-99"
        elif test_name == "Prothrombin Time":
            ref_range = "11.0-13.5"
        
        # Determine abnormal flag based on result value and reference range
        if test_name == "Urinalysis":
            abnormal_flag = "Normal" if result_value == "Negative" else "Abnormal"
        elif "<" in ref_range:
            threshold = float(ref_range.replace("<", ""))
            abnormal_flag = "Normal" if float(result_value) < threshold else "High"
        elif "-" in ref_range:
            low, high = map(float, ref_range.split("-"))
            if float(result_value) < low:
                abnormal_flag = "Low"
            elif float(result_value) > high:
                abnormal_flag = "High"
            else:
                abnormal_flag = "Normal"
        else:
            abnormal_flag = random.choice(ABNORMAL_FLAGS)
        
        # Sometimes use abbreviated flags
        if random.random() < 0.3:
            if abnormal_flag == "Normal":
                abnormal_flag = "N"
            elif abnormal_flag == "High":
                abnormal_flag = "H"
            elif abnormal_flag == "Low":
                abnormal_flag = "L"
        
        lab_result = {
            "patient_id": patient["patient_id"],
            "test_name": test_name,
            "test_date": generate_datetime(start_date, end_date),
            "result_value": result_value,
            "reference_range": ref_range,
            "unit": LAB_UNITS[test_name],
            "abnormal_flag": abnormal_flag,
            "ordering_provider": patient["provider_id"] if random.random() > 0.3 else generate_provider()[1]
        }
        
        # Introduce some data inconsistencies
        if random.random() < 0.1:  # 10% chance of different test name format
            lab_result["test_name"] = test_name.upper()
        if random.random() < 0.1 and lab_result["unit"]:  # 10% chance of missing unit
            lab_result["unit"] = ""
        
        lab_results.append(lab_result)
    
    return lab_results

def write_csv(data: List[Dict[str, Any]], filename: str):
    """Write data to a CSV file"""
    if not data:
        return
    
    with open(filename, 'w', newline='') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=data[0].keys())
        writer.writeheader()
        writer.writerows(data)
    
    print(f"Created {filename} with {len(data)} records")

def main():
    parser = argparse.ArgumentParser(description='Generate mock healthcare data')
    parser.add_argument('--patients', type=int, default=100, help='Number of patients to generate')
    parser.add_argument('--vitals', type=int, default=500, help='Number of vital sign records to generate')
    parser.add_argument('--medications', type=int, default=300, help='Number of medication records to generate')
    parser.add_argument('--labs', type=int, default=400, help='Number of lab result records to generate')
    parser.add_argument('--output', type=str, default='../sample_data', help='Output directory')
    args = parser.parse_args()
    
    # Create output directory if it doesn't exist
    os.makedirs(args.output, exist_ok=True)
    
    # Generate data
    print(f"Generating mock healthcare data...")
    patients = generate_patient_data(args.patients)
    vitals = generate_vitals_data(patients, args.vitals)
    medications = generate_medication_data(patients, args.medications)
    lab_results = generate_lab_data(patients, args.labs)
    
    # Write to CSV files
    write_csv(patients, os.path.join(args.output, 'patients.csv'))
    write_csv(vitals, os.path.join(args.output, 'vitals.csv'))
    write_csv(medications, os.path.join(args.output, 'medications.csv'))
    write_csv(lab_results, os.path.join(args.output, 'lab_results.csv'))
    
    print("Mock data generation complete!")

if __name__ == "__main__":
    main()
