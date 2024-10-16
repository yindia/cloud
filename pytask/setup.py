from setuptools import setup, find_packages


setup(name='taskcli',
      version='0.0.5',
      py_modules=['cloud','validate'],
      packages=find_packages(),  # Automatically find all packages under pytask
      include_package_data=True,
      install_requires=[
          'Click',
          'requests',
          'nose',
          'protobuf'
      ],
      package_data={'': ['*.txt', '*.lst']},
      entry_points={
          'console_scripts': [
              'taskcli=taskcli:cli',  # Ensure this points to the correct CLI
          ],
      },
      test_suite='nose.collector',
      tests_require=['nose'],
      )
