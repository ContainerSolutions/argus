resource "kubernetes_manifest" "network_control_1" {
  manifest = yamldecode(templatefile("requirement.yaml", 
  {
    description="Auto routing should be disabled",
    name="net-req-01"
    code="NET-REQ-01",
    version="1.0.0",
    componentclass="NetworkElement",
    }
    ))
}

resource "kubernetes_manifest" "network_control_2" {
  manifest = yamldecode(templatefile("requirement.yaml",
    {
    description="Outbound traffic should be sent to NSEC device",
    name="net-req-02"
    code="NET-REQ-02",
    version="1.0.0",
    componentclass="NetworkElement",
    }
))
}

resource "kubernetes_manifest" "network_control_3" {
  manifest = yamldecode(templatefile("requirement.yaml",
    {
    description="Public IP Addresses must not be used",
    name="net-req-03"
    code="NET-REQ-03",
    version="1.0.0",
    componentclass="NetworkElement",
    }
))
}

resource "kubernetes_manifest" "network_control_4" {
  manifest = yamldecode(templatefile("requirement.yaml",
    {
    description="Inbound internet facing communication must be behind central firewall",
    name="net-req-04"
    code="NET-REQ-04",
    version="1.0.0",
    componentclass="NetworkElement",
    }
))
}

resource "kubernetes_manifest" "network_control_5" {
  manifest = yamldecode(templatefile("requirement.yaml",
    {
    description="DNS requests must be microsegmented between network environments",
    name="net-req-05"
    code="NET-REQ-05",
    version="1.0.0",
    componentclass="NetworkElement",
    }
))
}

resource "kubernetes_manifest" "network_control_6" {
  manifest = yamldecode(templatefile("requirement.yaml",
    {
    name="net-req-06"
    code="NET-REQ-06",
    version="1.0.0",
    description="Any Inbound traffic needs to be properly managed behind firewall rules",
    componentclass="NetworkElement",
    }
))
}

resource "kubernetes_manifest" "network_control_7" {
  manifest = yamldecode(templatefile("requirement.yaml",
    {
    description="Networking design needs to ensure network isolation between different applications",
    name="net-req-07"
    code="NET-REQ-07",
    version="1.0.0",
    componentclass="NetworkElement",
    }
))
}

resource "kubernetes_manifest" "network_control_8" {
  manifest = yamldecode(templatefile("requirement.yaml",
    {
    description="Network traffic cannot leave the same bounded region/datacenter",
    name="net-req-08"
    code="NET-REQ-08",
    version="1.0.0",
    componentclass="NetworkElement",
    }
))
}

resource "kubernetes_manifest" "preventative_control" {
    count = 8
  manifest = yamldecode(templatefile("assessment.yaml",
    {
    name="req${count.index+1}-preventative",
    class= "PreventativeControl",
    code= "NET-REQ-0${count.index+1}",
    version="1.0.0",
    type="router"
    }
  ))
}

resource "kubernetes_manifest" "detective_control" {
    count = 8
  manifest = yamldecode(templatefile("assessment.yaml",
    {
    name="req${count.index+1}-detective",
    class= "DetectiveControl",
    code= "NET-REQ-0${count.index+1}",
    version="1.0.0",
    type="router"
    }
  ))
}

resource "kubernetes_manifest" "reactive_control" {
    count = 8
  manifest = yamldecode(templatefile("assessment.yaml",
    {
    name="req${count.index+1}-reactive",
    class= "ReactiveControl",
    code= "NET-REQ-0${count.index+1}",
    version="1.0.0",
    type="router"
    }
  ))
}

resource "kubernetes_manifest" "reasoning_attestation" {
    count = 8
  manifest = yamldecode(templatefile("attestation.yaml",
    {
    name="req${count.index+1}-detective",
    implementation = "req${count.index+1}-detective",
    }
  ))
}

resource "kubernetes_manifest" "deployed_attestation" {
    count = 8
  manifest = yamldecode(templatefile("attestation.yaml",
    {
    name="req${count.index+1}-preventative",
    implementation = "req${count.index+1}-preventative",
    }
  ))
}

resource "kubernetes_manifest" "reactive_reasoning_attestation" {
    count = 8
  manifest = yamldecode(templatefile("attestation.yaml",
    {
    name="req${count.index+1}-reactive-reasoning",
    implementation = "req${count.index+1}-reactive",
    }
  ))
}

resource "kubernetes_manifest" "reactive_deployed_attestation" {
    count = 8
  manifest = yamldecode(templatefile("attestation.yaml",
    {
    name="req${count.index+1}-reactive-deployed",
    implementation = "req${count.index+1}-reactive",
    }
  ))
}

resource "kubernetes_manifest" "detective_reasoning_attestation" {
    count = 8
  manifest = yamldecode(templatefile("attestation.yaml",
    {
    name="req${count.index+1}-detective-reasoning",
    implementation = "req${count.index+1}-detective",
    }
  ))
}

resource "kubernetes_manifest" "detective_deployed_attestation" {
    count = 8
  manifest = yamldecode(templatefile("attestation.yaml",
    {
    name="req${count.index+1}-detective-deployed",
    implementation = "req${count.index+1}-detective",
    }
  ))
}
