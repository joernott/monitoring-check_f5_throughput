#
# spec file for package f5_throughput
#
Name:           icinga2-plugin-check_f5_throughput
Version:        %{version}
Release:        %{release}
Summary:        Monitoring f5 throughput with Icinga2/Nagios
License:        BSD
Group:          Sytem/Utilities
Vendor:         Ott-Consult UG
Packager:       Joern Ott
Url:            https://github.com/joernott/check_f5_throughput
Source0:        check_f5_throughput-%{version}.tar.gz
BuildArch:      x86_64

%description
A check for Icinga2 or Nagios to monitor the throughput of an F5 BigIP
loadbalancer using SNMP.

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/usr/lib64/nagios/plugins /etc/icinga2/
cd $RPM_BUILD_ROOT
tar -xzf %{SOURCE0}

%files
%defattr(-,root,root,755)
%attr(755, root, root) /usr/lib64/nagios/plugins/check_f5_throughput
%config(noreplace) /etc/icinga2/check_f5_throughput.yaml

%changelog
* Thu Feb 18 2021 Joern Ott <joern.ott@ott-consult.de>
- Initial version