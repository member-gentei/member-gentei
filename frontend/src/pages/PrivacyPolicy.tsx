import React, { Fragment } from "react";
import Footer from "../components/Footer";
import { Box, Container, Link, Sheet, Typography } from "@mui/joy";

function PrivacyPolicy() {
  return (
    <Fragment>
      <Sheet
        variant="solid"
        sx={{
          backgroundColor: "rgb(39, 128, 227)",
          padding: 4,
          mb: 12,
        }}
      >
        <Box sx={{ mt: 8 }}>
          <Typography level="h1" sx={{ color: "white" }}>
            Gentei / 限定
          </Typography>
          <Typography level="title-lg" sx={{ color: "white" }}>
            VTuber channel membership verification
          </Typography>
        </Box>
      </Sheet>
      <Container component="section" role="main">
        <section role="main" className="section">
          <div className="container">
            <div className="content">
              <h1 className="title">Privacy Policy</h1>
              <h2 className="title">In colloquial form</h2>
              <p>
                We collect the following personally identifiable information
                (PII) to provide and improve this service.
              </p>
              <ul>
                <li>
                  Platform user identifier (e.g. Discord user snowflake IDs and
                  YouTube channel IDs)
                </li>
                <li>Usernames and channel names</li>
              </ul>
              <p>
                The service requires transmission of information to Google and
                Discord to verify membership status, which you can revoke at any
                time via either requesting account deletion or using the{" "}
                <Link href="https://myaccount.google.com/permissions">
                  Google Security Settings
                </Link>{" "}
                page and Discord Authorized Apps pages respectively.
              </p>
              <p>Yeah, that's it. No sale of information or anything.</p>
              <hr />
              <h2 className="title">In formal form</h2>
              <div className="legalese">
                <p>
                  ignacio.io LLC built the Member Gentei app ("Gentei") as a
                  Free app. This SERVICE is provided by ignacio.io LLC at no
                  cost and is intended for use as is.
                </p>
                <p>
                  This page is used to inform visitors regarding our policies
                  with the collection, use, and disclosure of Personal
                  Information if anyone decided to use our Service.
                </p>
                <p>
                  If you choose to use our Service, then you agree to the
                  collection and use of information in relation to this policy.
                  The Personal Information that we collect is used for providing
                  and improving the Service. We will not use or share your
                  information with anyone except as described in this Privacy
                  Policy.
                </p>
                <p>
                  The terms used in this Privacy Policy have the same meanings
                  as in our Terms and Conditions, which is accessible at Member
                  Gentei unless otherwise defined in this Privacy Policy.
                </p>
                <p>
                  <strong>Information Collection and Use</strong>
                </p>
                <p>
                  For a better experience, while using our Service, we may
                  require you to provide us with certain personally identifiable
                  information, including but not limited to Platform user
                  identifiers. The information that we request will be retained
                  by us and used as described in this privacy policy.
                </p>
                <p>
                  <strong>Log Data</strong>
                </p>
                <p>
                  We want to inform you that whenever you use our Service, in a
                  case of an error in the app we collect data and information
                  (through third party products) on your phone called Log Data.
                  This Log Data may include information such as your device
                  Internet Protocol (“IP”) address, device name, operating
                  system version, the configuration of the app when utilizing
                  our Service, the time and date of your use of the Service, and
                  other statistics.
                </p>
                <p>
                  <strong>Cookies</strong>
                </p>
                <p>
                  Cookies are files with a small amount of data that are
                  commonly used as anonymous unique identifiers. These are sent
                  to your browser from the websites that you visit and are
                  stored on your device's internal memory.
                </p>
                <p>
                  This Service does not use these “cookies” explicitly. However,
                  the app may use third party code and libraries that use
                  “cookies” to collect information and improve their services.
                  You have the option to either accept or refuse these cookies
                  and know when a cookie is being sent to your device. If you
                  choose to refuse our cookies, you may not be able to use some
                  portions of this Service.
                </p>
                <p>
                  <strong>Service Providers</strong>
                </p>
                <p>
                  We may employ third-party companies and individuals due to the
                  following reasons:
                </p>
                <ul>
                  <li>To facilitate our Service;</li>
                  <li>To provide the Service on our behalf;</li>
                  <li>To perform Service-related services; or</li>
                  <li>To assist us in analyzing how our Service is used.</li>
                </ul>
                <p>
                  We want to inform users of this Service that these third
                  parties have access to your Personal Information. The reason
                  is to perform the tasks assigned to them on our behalf.
                  However, they are obligated not to disclose or use the
                  information for any other purpose.
                </p>
                <p>
                  Gentei uses YouTube API Services in order to provide the
                  Service. Users that connect their account to YouTube agree to
                  be bound by YouTube's Terms of Service, which can be found at{" "}
                  <Link href="https://www.youtube.com/t/terms">
                    https://www.youtube.com/t/terms
                  </Link>
                  . This includes Google's Privacy Policy, available at{" "}
                  <Link href="https://policies.google.com/privacy">
                    https://policies.google.com/privacy
                  </Link>
                  .
                </p>
                <p>
                  <strong>Security</strong>
                </p>
                <p>
                  We value your trust in providing us your Personal Information,
                  thus we are striving to use commercially acceptable means of
                  protecting it. But remember that no method of transmission
                  over the internet, or method of electronic storage is 100%
                  secure and reliable, and we cannot guarantee its absolute
                  security.
                </p>
                <p>
                  <strong>Links to Other Sites</strong>
                </p>
                <p>
                  This Service may contain links to other sites. If you click on
                  a third-party link, you will be directed to that site. Note
                  that these external sites are not operated by us. Therefore,
                  we strongly advise you to review the Privacy Policy of these
                  websites. We have no control over and assume no responsibility
                  for the content, privacy policies, or practices of any
                  third-party sites or services.
                </p>
                <p>
                  <strong>Children's Privacy</strong>
                </p>
                <p>
                  These Services do not address anyone under the age of 13. We
                  do not knowingly collect personally identifiable information
                  from children under 13. In the case we discover that a child
                  under 13 has provided us with personal information, we
                  immediately delete this from our servers. If you are a parent
                  or guardian and you are aware that your child has provided us
                  with personal information, please contact us so that we will
                  be able to do necessary actions.
                </p>
                <p>
                  <strong>Changes to This Privacy Policy</strong>
                </p>
                <p>
                  We may update our Privacy Policy from time to time. Thus, you
                  are advised to review this page periodically for any changes.
                  We will notify you of any changes by posting the new Privacy
                  Policy on this page.
                </p>
                <p>This policy is effective as of 2024-09-03.</p>
                <p>
                  <strong>Contact Us</strong>
                </p>
                <p>
                  If you have any questions or suggestions about our Privacy
                  Policy, do not hesitate to contact us at
                  member-gentei@ignacio.io.
                </p>
              </div>
            </div>
          </div>
        </section>
        <Footer />
      </Container>
    </Fragment>
  );
}

export default PrivacyPolicy;
