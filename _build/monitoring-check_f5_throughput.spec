#
# spec file for package monitoring-check_f5_throughput
#
Name:           monitoring-check_f5_throughput
Version:        %{version}
Release:        %{release}
Summary:        Monitoring f5 throughput with Icinga2/Nagios
License:        BSD
Group:          Sytem/Utilities
Vendor:         Ott-Consult UG
Packager:       Joern Ott
Url:            https://github.com/joernott/monitoring-check_f5_throughput
Source0:        monitoring-check_f5_throughput-%{version}.tar.gz
BuildArch:      x86_64

%description
A check for Icinga2 or Nagios to monitor the throughput of an F5 BigIP
loadbalancer using SNMP.

%prep
cd "$RPM_BUILD_DIR"
rm -rf *
tar -xzf "%{SOURCE0}"
STATUS=$?
if [ $STATUS -ne 0 ]; then
  exit $STATUS
fi
/usr/bin/chmod -Rf a+rX,u+w,g-w,o-w .

%build
cd "$RPM_BUILD_DIR"
go build -v -o %{name}

%install
install -Dpm 0755 check_f5_throughput "%{buildroot}/usr/lib64/nagios/plugins/check_f5_throughput"
install -Dpm 0640 sample_config.yaml "%{buildroot}/etc/icinga2/check_f5_throughput.yaml"

%files
%defattr(-,root,root,755)
%attr(755, root, root) /usr/lib64/nagios/plugins/check_f5_throughput
%config(noreplace) /etc/icinga2/check_f5_throughput.yaml

%changelog
* Thu Apr 28 2022 Joern Ott <joern.ott@ott-consult.de>
- Standardize rpm build process and rename repo
* Thu Feb 18 2021 Joern Ott <joern.ott@ott-consult.de>
- Initial version
